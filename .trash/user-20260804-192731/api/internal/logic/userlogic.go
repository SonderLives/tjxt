package logic

import (
	"context"

	"tjxt/pkg/response"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 学员注册 ============

type StudentRegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudentRegisterLogic {
	return &StudentRegisterLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *StudentRegisterLogic) StudentRegister(req *types.StudentFormDTO) (resp *result.R, err error) {
	if err := l.svcCtx.UserService.Register(l.ctx, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 修改学员密码 ============

type StudentPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudentPasswordLogic {
	return &StudentPasswordLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *StudentPasswordLogic) StudentPassword(req *types.StudentFormDTO) (resp *result.R, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserService.ChangePassword(l.ctx, userID, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 学员分页 ============

type StudentPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStudentPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StudentPageLogic {
	return &StudentPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *StudentPageLogic) StudentPage(req *types.UserPageReq) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	page, err := l.svcCtx.UserService.PageUsers(l.ctx, 2, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 教师分页 ============

type TeacherPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTeacherPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TeacherPageLogic {
	return &TeacherPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *TeacherPageLogic) TeacherPage(req *types.UserPageReq) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	page, err := l.svcCtx.UserService.PageUsers(l.ctx, 3, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 员工分页 ============

type StaffPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewStaffPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StaffPageLogic {
	return &StaffPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *StaffPageLogic) StaffPage(req *types.UserPageReq) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	page, err := l.svcCtx.UserService.PageUsers(l.ctx, 1, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 当前用户信息 ============

type UserMeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserMeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserMeLogic {
	return &UserMeLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserMeLogic) UserMe() (resp *result.R, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	vo, err := l.svcCtx.UserService.GetMe(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 新增用户 ============

type UserCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCreateLogic {
	return &UserCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserCreateLogic) UserCreate(req *types.UserDTO) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	operatorID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	id, err := l.svcCtx.UserService.CreateUser(l.ctx, operatorID, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(id), nil
}

// ============ 更新当前用户 ============

type UserUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserUpdateLogic {
	return &UserUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserUpdateLogic) UserUpdate(req *types.UserFormDTO) (resp *result.R, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserService.UpdateMe(l.ctx, userID, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 校验手机号 ============

type UserCheckCellphoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCheckCellphoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCheckCellphoneLogic {
	return &UserCheckCellphoneLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserCheckCellphoneLogic) UserCheckCellphone(req *types.UserCheckCellphoneReq) (resp *result.R, err error) {
	exists, err := l.svcCtx.UserService.CheckCellphone(l.ctx, req.Cellphone)
	if err != nil {
		return nil, err
	}
	return result.OkData(exists), nil
}

// ============ 按 id 查询用户 ============

type UserGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserGetLogic {
	return &UserGetLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserGetLogic) UserGet(req *types.UserIdReq) (resp *result.R, err error) {
	user, err := l.svcCtx.UserService.GetUser(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(user), nil
}

// ============ 管理端更新用户 ============

type UserAdminUpdateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserAdminUpdateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserAdminUpdateLogic {
	return &UserAdminUpdateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserAdminUpdateLogic) UserAdminUpdate(req *types.UserIdReq, body *types.UserDTO) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserService.AdminUpdate(l.ctx, req.Id, body); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 重置密码 ============

type UserResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserResetPasswordLogic {
	return &UserResetPasswordLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserResetPasswordLogic) UserResetPassword(req *types.UserIdReq) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserService.ResetPassword(l.ctx, req.Id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 修改用户状态 ============

type UserStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserStatusLogic {
	return &UserStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *UserStatusLogic) UserStatus(req *types.UserStatusReq) (resp *result.R, err error) {
	if err := requireStaffOrAdmin(l.ctx); err != nil {
		return nil, err
	}
	if err := l.svcCtx.UserService.UpdateStatus(l.ctx, req.Id, req.Status); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}
