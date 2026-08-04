package service

import (
	"context"
	"database/sql"
	"time"

	"common/idgen"
	"common/xerr"
	"user/internal/model"
	"user/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"golang.org/x/crypto/bcrypt"
)

// SortByWhitelist 允许的排序字段
var SortByWhitelist = map[string]struct{}{
	"id":          {},
	"create_time": {},
}

// ValidateSortBy 校验排序字段是否合法
func ValidateSortBy(sortBy string) error {
	if sortBy == "" {
		return nil
	}
	if _, ok := SortByWhitelist[sortBy]; !ok {
		return xerr.BadRequestf("sortBy 非法: %s", sortBy)
	}
	return nil
}

// UserService 用户业务接口
type UserService interface {
	// Register 学员注册
	Register(ctx context.Context, req *types.StudentFormDTO) error
	// ChangePassword 修改学员密码（校验旧密码）
	ChangePassword(ctx context.Context, userId int64, req *types.StudentFormDTO) error
	// GetMe 获取当前登录用户信息
	GetMe(ctx context.Context, userId int64) (*types.UserDetailVO, error)
	// GetUser 按 id 查询用户（供内部服务/管理端调用）
	GetUser(ctx context.Context, id int64) (*types.UserDTO, error)
	// GetUserList 批量查询用户详情
	GetUserList(ctx context.Context, ids []int64) (map[int64]*types.UserDTO, error)
	// CheckCellphone 校验手机号是否已注册
	CheckCellphone(ctx context.Context, cellPhone string) (bool, error)
	// UpdateMe 更新当前登录用户资料（可改密码）
	UpdateMe(ctx context.Context, userId int64, req *types.UserFormDTO) error
	// CreateUser 新增员工/教师（事务）
	CreateUser(ctx context.Context, operatorId int64, req *types.UserDTO) (int64, error)
	// AdminUpdate 管理端更新用户资料
	AdminUpdate(ctx context.Context, id int64, req *types.UserDTO) error
	// ResetPassword 重置密码为默认密码
	ResetPassword(ctx context.Context, id int64) error
	// UpdateStatus 修改用户状态
	UpdateStatus(ctx context.Context, id, status int64) error
	// PageUsers 分页查询指定类型用户
	PageUsers(ctx context.Context, userType int64, req *types.UserPageReq) (*types.Page, error)
}

type userService struct {
	userModel       *model.UserModel
	detailModel     *model.UserDetailModel
	authModel       *model.AuthModel
	defaultPassword string
}

// NewUserService 创建用户业务服务
func NewUserService(userModel *model.UserModel, detailModel *model.UserDetailModel, authModel *model.AuthModel, defaultPassword string) UserService {
	return &userService{
		userModel:       userModel,
		detailModel:     detailModel,
		authModel:       authModel,
		defaultPassword: defaultPassword,
	}
}

// Register 学员注册（事务）
func (s *userService) Register(ctx context.Context, req *types.StudentFormDTO) error {
	if req.CellPhone == "" || req.Password == "" {
		return xerr.BadRequestf("手机号与密码不能为空")
	}
	existing, err := s.userModel.FindByCellPhone(ctx, req.CellPhone)
	if err == nil && existing != nil {
		return xerr.Conflict("该手机号已注册")
	}
	if err != nil && err != sql.ErrNoRows {
		return xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "密码加密失败")
	}

	now := time.Now()
	id := idgen.NextID()
	return s.userModel.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		if err := s.userModel.InsertTx(ctx, tx, &model.User{
			Id:         id,
			Username:   req.CellPhone,
			CellPhone:  req.CellPhone,
			Password:   string(hash),
			Type:       model.UserTypeStudent,
			Status:     model.UserStatusNormal,
			CreateTime: now,
		}); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "注册失败")
		}
		if err := s.detailModel.InsertTx(ctx, tx, &model.UserDetail{
			Id:           id,
			Type:         model.UserTypeStudent,
			Name:         maskPhone(req.CellPhone),
			Gender:       0,
			RoleId:       studentRoleID,
			CourseAmount: 0,
			CreateTime:   now,
		}); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "注册失败")
		}
		return nil
	})
}

// ChangePassword 修改学员密码
func (s *userService) ChangePassword(ctx context.Context, userId int64, req *types.StudentFormDTO) error {
	_, err := s.userModel.FindById(ctx, userId)
	if err == sql.ErrNoRows {
		return xerr.NotFound("账户不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "密码加密失败")
	}
	return s.userModel.UpdatePassword(ctx, userId, string(hash))
}

// GetMe 获取当前登录用户信息
func (s *userService) GetMe(ctx context.Context, userId int64) (*types.UserDetailVO, error) {
	u, err := s.userModel.FindById(ctx, userId)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("账户不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	d, err := s.detailModel.FindById(ctx, userId)
	if err != nil && err != sql.ErrNoRows {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户详情失败")
	}
	vo := &types.UserDetailVO{
		Id:         u.Id,
		Username:   u.Username,
		CellPhone:  u.CellPhone,
		CreateTime: u.CreateTime.Format(time.RFC3339),
	}
	if d != nil {
		vo.Name = d.Name
		vo.Gender = d.Gender
		vo.Icon = d.Icon
		vo.Email = d.Email
		vo.Qq = d.Qq
		vo.Province = d.Province
		vo.City = d.City
		vo.District = d.District
		vo.Intro = d.Intro
		roleName, _ := s.authModel.RoleName(ctx, d.RoleId)
		vo.RoleName = roleName
	}
	return vo, nil
}

// GetUser 按 id 查询用户（内部服务返回 UserDTO）
func (s *userService) GetUser(ctx context.Context, id int64) (*types.UserDTO, error) {
	u, err := s.userModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("用户不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	dto := &types.UserDTO{
		Id:        u.Id,
		Username:  u.Username,
		CellPhone: u.CellPhone,
		Type:      u.Type,
	}
	d, err := s.detailModel.FindById(ctx, id)
	if err == nil && d != nil {
		dto.Name = d.Name
		dto.Gender = d.Gender
		dto.Icon = d.Icon
		dto.Email = d.Email
		dto.Qq = d.Qq
		dto.Job = d.Job
		dto.Province = d.Province
		dto.City = d.City
		dto.District = d.District
		dto.Intro = d.Intro
		dto.Photo = d.Photo
		dto.RoleId = d.RoleId
	}
	return dto, nil
}

// GetUserList 批量查询用户详情
func (s *userService) GetUserList(ctx context.Context, ids []int64) (map[int64]*types.UserDTO, error) {
	if len(ids) == 0 {
		return make(map[int64]*types.UserDTO), nil
	}

	users, err := s.userModel.FindListByIds(ctx, ids)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}

	details, err := s.detailModel.ListByIDs(ctx, ids)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询详情失败")
	}

	result := make(map[int64]*types.UserDTO, len(users))
	for _, u := range users {
		dto := &types.UserDTO{
			Id:        u.Id,
			Username:  u.Username,
			CellPhone: u.CellPhone,
			Type:      u.Type,
		}
		if d, ok := details[u.Id]; ok {
			dto.Name = d.Name
			dto.Gender = d.Gender
			dto.Icon = d.Icon
			dto.Email = d.Email
			dto.Qq = d.Qq
			dto.Job = d.Job
			dto.Province = d.Province
			dto.City = d.City
			dto.District = d.District
			dto.Intro = d.Intro
			dto.Photo = d.Photo
			dto.RoleId = d.RoleId
		}
		result[u.Id] = dto
	}
	return result, nil
}

// CheckCellphone 校验手机号是否已注册
func (s *userService) CheckCellphone(ctx context.Context, cellPhone string) (bool, error) {
	u, err := s.userModel.FindByCellPhone(ctx, cellPhone)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	return u != nil, nil
}

// UpdateMe 更新当前登录用户资料
func (s *userService) UpdateMe(ctx context.Context, userId int64, req *types.UserFormDTO) error {
	u, err := s.userModel.FindById(ctx, userId)
	if err == sql.ErrNoRows {
		return xerr.NotFound("账户不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	// 改密码：校验旧密码
	if req.Password != "" {
		if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.OldPassword)) != nil {
			return xerr.BadRequestf("原密码错误")
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "密码加密失败")
		}
		if err := s.userModel.UpdatePassword(ctx, userId, string(hash)); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "修改密码失败")
		}
	}
	// 更新资料
	d, err := s.detailModel.FindById(ctx, userId)
	if err == sql.ErrNoRows {
		d = &model.UserDetail{Id: userId, Type: u.Type}
	}
	applyDetail(d, req.UserDTO)
	if err := s.detailModel.UpdateDetail(ctx, d); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新资料失败")
	}
	return nil
}

// CreateUser 新增员工/教师（事务）
func (s *userService) CreateUser(ctx context.Context, operatorId int64, req *types.UserDTO) (int64, error) {
	if req.CellPhone == "" || req.Username == "" {
		return 0, xerr.BadRequestf("手机号与用户名不能为空")
	}
	if req.Type != model.UserTypeStaff && req.Type != model.UserTypeTeacher {
		return 0, xerr.BadRequestf("仅支持新增员工或教师")
	}
	existing, err := s.userModel.FindByCellPhone(ctx, req.CellPhone)
	if err == nil && existing != nil {
		return 0, xerr.Conflict("该手机号已存在")
	}
	if err != nil && err != sql.ErrNoRows {
		return 0, xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(s.defaultPwd()), bcrypt.DefaultCost)
	if err != nil {
		return 0, xerr.Wrap(err, xerr.CodeInternal, "密码加密失败")
	}
	now := time.Now()
	id := idgen.NextID()

	return id, s.userModel.TransactCtx(ctx, func(ctx context.Context, tx sqlx.Session) error {
		if err := s.userModel.InsertTx(ctx, tx, &model.User{
			Id:         id,
			Username:   req.Username,
			CellPhone:  req.CellPhone,
			Password:   string(hash),
			Type:       req.Type,
			Status:     model.UserStatusNormal,
			CreateTime: now,
			Creater:    sql.NullInt64{Int64: operatorId, Valid: operatorId > 0},
			Updater:    sql.NullInt64{Int64: operatorId, Valid: operatorId > 0},
		}); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "新增用户失败")
		}
		detail := &model.UserDetail{
			Id:           id,
			Type:         req.Type,
			Name:         firstNonEmpty(req.Name, req.Username),
			Gender:       req.Gender,
			Icon:         req.Icon,
			Email:        req.Email,
			Qq:           req.Qq,
			Job:          req.Job,
			Province:     req.Province,
			City:         req.City,
			District:     req.District,
			Intro:        req.Intro,
			Photo:        req.Photo,
			RoleId:       firstPositive(req.RoleId, defaultRoleID(req.Type)),
			CourseAmount: 0,
			CreateTime:   now,
		}
		if err := s.detailModel.InsertTx(ctx, tx, detail); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "新增用户失败")
		}
		return nil
	})
}

// AdminUpdate 管理端更新用户资料
func (s *userService) AdminUpdate(ctx context.Context, id int64, req *types.UserDTO) error {
	u, err := s.userModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("用户不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	if err := s.userModel.UpdateAccount(ctx, id, req.Username, req.CellPhone); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新账户失败")
	}
	d, err := s.detailModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		d = &model.UserDetail{Id: id, Type: u.Type}
	}
	applyDetail(d, *req)
	if err := s.detailModel.UpdateDetail(ctx, d); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新用户失败")
	}
	return nil
}

// ResetPassword 重置密码为默认密码
func (s *userService) ResetPassword(ctx context.Context, id int64) error {
	if _, err := s.userModel.FindById(ctx, id); err != nil {
		if err == sql.ErrNoRows {
			return xerr.NotFound("用户不存在")
		}
		return xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(s.defaultPwd()), bcrypt.DefaultCost)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "密码加密失败")
	}
	return s.userModel.UpdatePassword(ctx, id, string(hash))
}

// UpdateStatus 修改用户状态
func (s *userService) UpdateStatus(ctx context.Context, id, status int64) error {
	if status != model.UserStatusNormal && status != model.UserStatusDisabled {
		return xerr.BadRequestf("非法的状态值")
	}
	if err := s.userModel.UpdateStatus(ctx, id, status); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "修改状态失败")
	}
	return nil
}

// PageUsers 分页查询指定类型用户
func (s *userService) PageUsers(ctx context.Context, userType int64, req *types.UserPageReq) (*types.Page, error) {
	// 校验 SortBy
	if err := ValidateSortBy(req.SortBy); err != nil {
		return nil, err
	}

	offset, limit := normalizePage(req.PageNo, req.PageSize)
	cond := &model.PageCond{
		UserType: userType,
		Name:     req.Name,
		Phone:    req.Phone,
		Status:   req.Status,
		Offset:   offset,
		Limit:    limit,
		IsAsc:    req.IsAsc,
		SortBy:   req.SortBy,
	}
	rows, total, err := s.detailModel.ListPage(ctx, cond)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户失败")
	}

	switch userType {
	case model.UserTypeStudent:
		list := make([]types.StudentPageVO, 0, len(rows))
		for i := range rows {
			r := &rows[i]
			list = append(list, types.StudentPageVO{
				Id:           r.Id,
				Name:         r.Name,
				CellPhone:    r.CellPhone,
				Gender:       r.Gender,
				Icon:         r.Icon,
				Status:       r.Status,
				CourseAmount: r.CourseAmount,
				CreateTime:   r.CreateTime.Format(time.RFC3339),
			})
		}
		return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
	case model.UserTypeTeacher:
		list := make([]types.TeacherPageVO, 0, len(rows))
		for i := range rows {
			r := &rows[i]
			list = append(list, types.TeacherPageVO{
				Id:           r.Id,
				Name:         r.Name,
				CellPhone:    r.CellPhone,
				Icon:         r.Icon,
				Photo:        r.Photo,
				Job:          r.Job,
				Intro:        r.Intro,
				Status:       r.Status,
				CourseAmount: r.CourseAmount,
				CreateTime:   r.CreateTime.Format(time.RFC3339),
			})
		}
		return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
	default: // 员工
		list := make([]types.StaffVO, 0, len(rows))
		for i := range rows {
			r := &rows[i]
			roleName, _ := s.authModel.RoleName(ctx, r.RoleId)
			logx.Debugf("RoleName lookup: roleId=%d roleName=%q", r.RoleId, roleName)
			list = append(list, types.StaffVO{
				Id:         r.Id,
				Name:       r.Name,
				CellPhone:  r.CellPhone,
				Icon:       r.Icon,
				RoleId:     r.RoleId,
				RoleName:   roleName,
				Status:     r.Status,
				CreateTime: r.CreateTime.Format(time.RFC3339),
			})
		}
		return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
	}
}

func (s *userService) defaultPwd() string {
	if s.defaultPassword == "" {
		return "123456"
	}
	return s.defaultPassword
}

// ==================== 公共 helper ====================

// normalizePage 归一化分页参数
func normalizePage(pageNo, pageSize int64) (offset, limit int64) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return (pageNo - 1) * pageSize, pageSize
}

func calcPages(total, pageSize int64) int64 {
	if pageSize <= 0 || total <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	return pages
}

func applyDetail(d *model.UserDetail, dto types.UserDTO) {
	if dto.Name != "" {
		d.Name = dto.Name
	}
	if dto.RoleId > 0 {
		d.RoleId = dto.RoleId
	}
	d.Gender = dto.Gender
	d.Icon = dto.Icon
	d.Email = dto.Email
	d.Qq = dto.Qq
	d.Job = dto.Job
	d.Province = dto.Province
	d.City = dto.City
	d.District = dto.District
	d.Intro = dto.Intro
	d.Photo = dto.Photo
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func firstPositive(a, b int64) int64 {
	if a > 0 {
		return a
	}
	return b
}

const studentRoleID int64 = 2

func defaultRoleID(userType int64) int64 {
	switch userType {
	case model.UserTypeStaff:
		return 1
	case model.UserTypeTeacher:
		return 3
	default:
		return studentRoleID
	}
}

func maskPhone(phone string) string {
	if len(phone) != 11 {
		return phone
	}
	return phone[:3] + "****" + phone[7:]
}