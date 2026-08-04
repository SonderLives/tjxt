package service

import (
	"context"
	"database/sql"
	"sort"
	"strings"

	"common/idgen"
	"common/page"
	"common/result"
	"common/xerr"
	"course/internal/model"
	"course/internal/types"
)

// CourseService 课程业务接口（正式表 + 草稿表，与 Java 端 ICourseService/ICourseDraftService 对齐）。
type CourseService interface {
	// SaveBaseInfo 保存课程基本信息（新增或编辑草稿）。
	SaveBaseInfo(ctx context.Context, req *types.CourseBaseInfoSaveDTO, userId int64) (*types.CourseSaveVO, error)
	// GetBaseInfo 获取课程基础信息。see=true 优先读取正式数据（查看页），否则读草稿（编辑页）。
	GetBaseInfo(ctx context.Context, id int64, see bool) (*types.CourseBaseInfoVO, error)
	// CheckName 校验课程名称是否已存在。
	CheckName(ctx context.Context, name string, id int64) (*types.NameExistVO, error)
	// Delete 删除课程（正式数据 + 草稿 + 关联数据）。
	Delete(ctx context.Context, id int64) error
	// QueryForPage 管理端课程分页搜索。status=1/3 查草稿，status=2/4 查正式数据。
	QueryForPage(ctx context.Context, req *types.CoursePageQuery) (*result.Page, error)
	// QueryCourseIdByName 按名称模糊查询课程 id 列表。
	QueryCourseIdByName(ctx context.Context, name string) ([]int64, error)
	// QuerySimpleInfoList 按条件列表查询课程简要信息。
	QuerySimpleInfoList(ctx context.Context, ids, thirdCataIds []int64) ([]*types.CourseSimpleInfoDTO, error)

	// ---------- 上架/下架 ----------
	// CheckBeforeUpShelf 课程上架前校验。
	CheckBeforeUpShelf(ctx context.Context, id int64) error
	// UpShelf 课程上架（草稿拷贝到正式表）。
	UpShelf(ctx context.Context, id int64) error
	// DownShelf 课程下架（正式数据拷贝回草稿）。
	DownShelf(ctx context.Context, id int64) error

	// ---------- 目录/媒资/题目 ----------
	// SaveCatas 保存课程章节目录。
	SaveCatas(ctx context.Context, id int64, step int64, list []*types.CataSaveDTO) error
	// QueryCatas 获取课程章节（see=true 优先正式，否则草稿）。
	QueryCatas(ctx context.Context, id int64, see, withPractice bool) ([]*types.CataVO, error)
	// SaveMedia 保存小节媒资信息。
	SaveMedia(ctx context.Context, id int64, list []*types.CourseMediaSaveDTO) error
	// SaveSubjects 保存小节/练习的题目。
	SaveSubjects(ctx context.Context, id int64, list []*types.CataSubjectDTO) error
	// GetSubjects 获取小节/练习中的题目。
	GetSubjects(ctx context.Context, id int64) ([]*types.CataSimpleSubjectVO, error)

	// ---------- 老师 ----------
	// SaveTeachers 保存课程老师关系。
	SaveTeachers(ctx context.Context, req *types.CourseTeacherSaveDTO, userId int64) error
	// QueryTeachers 查询课程相关老师信息。
	QueryTeachers(ctx context.Context, id int64, see bool) ([]*types.CourseTeacherVO, error)

	// ---------- 内部调用 ----------
	// GetInfoById 获取课程详细信息（含目录、老师）。
	GetInfoById(ctx context.Context, id int64, withCatalogue, withTeachers bool) (*types.CourseFullInfoDTO, error)
	// GetCourseDTOById 获取课程信息（搜索索引库使用）。
	GetCourseDTOById(ctx context.Context, id int64) (*types.CourseDTO, error)
	// QueryCourseAndCatalog 查询课程基本信息、目录、学习进度。
	QueryCourseAndCatalog(ctx context.Context, id int64) (*types.CourseAndSectionVO, error)
	// CountSubjectNumAndCourseNumOfTeacher 统计老师名下课程数与出题数。
	CountSubjectNumAndCourseNumOfTeacher(ctx context.Context, teacherIds []int64) ([]*types.SubNumAndCourseNumDTO, error)
	// UpdateStep 更新课程编辑进度（内部使用）。
	UpdateStep(ctx context.Context, id int64, step int64) error
}

type courseService struct {
	courseDraftModel         *model.CourseDraftModel
	courseModel              *model.CourseModel
	courseContentDraftModel  *model.CourseContentDraftModel
	courseContentModel       *model.CourseContentModel
	cataDraftModel           *model.CourseCatalogueDraftModel
	cataModel                *model.CourseCatalogueModel
	teacherDraftModel        *model.CourseTeacherDraftModel
	teacherModel             *model.CourseTeacherModel
	cataSubjectDraftModel    *model.CourseCataSubjectDraftModel
	categoryModel            *model.CategoryModel
	subjectModel             *model.SubjectModel
	userDetailModel          *model.UserDetailModel
	trailerDuration          int64
}

// NewCourseService 创建课程业务服务。
func NewCourseService(
	courseDraftModel *model.CourseDraftModel,
	courseModel *model.CourseModel,
	courseContentDraftModel *model.CourseContentDraftModel,
	courseContentModel *model.CourseContentModel,
	cataDraftModel *model.CourseCatalogueDraftModel,
	cataModel *model.CourseCatalogueModel,
	teacherDraftModel *model.CourseTeacherDraftModel,
	teacherModel *model.CourseTeacherModel,
	cataSubjectDraftModel *model.CourseCataSubjectDraftModel,
	categoryModel *model.CategoryModel,
	subjectModel *model.SubjectModel,
	userDetailModel *model.UserDetailModel,
	trailerDuration int64,
) CourseService {
	return &courseService{
		courseDraftModel:        courseDraftModel,
		courseModel:             courseModel,
		courseContentDraftModel: courseContentDraftModel,
		courseContentModel:      courseContentModel,
		cataDraftModel:          cataDraftModel,
		cataModel:               cataModel,
		teacherDraftModel:       teacherDraftModel,
		teacherModel:            teacherModel,
		cataSubjectDraftModel:   cataSubjectDraftModel,
		categoryModel:           categoryModel,
		subjectModel:            subjectModel,
		userDetailModel:         userDetailModel,
		trailerDuration:         trailerDuration,
	}
}

// SaveBaseInfo 保存课程基本信息。
func (s *courseService) SaveBaseInfo(ctx context.Context, req *types.CourseBaseInfoSaveDTO, userId int64) (*types.CourseSaveVO, error) {
	var course *model.Course
	categoryIdList := make([]int64, 0, 3)
	if req.Id > 0 {
		c, err := s.courseModel.FindById(ctx, req.Id)
		if err != nil && err != sql.ErrNoRows {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
		}
		course = c
	}
	if req.Id == 0 || course == nil {
		// 新增或编辑不存在的课程：校验同名、校验分类
		if err := s.checkSameName(ctx, req.Id, req.Name); err != nil {
			return nil, err
		}
		ids, err := s.checkCategory(ctx, req.ThirdCateId)
		if err != nil {
			return nil, err
		}
		categoryIdList = ids
	}

	draft := &model.CourseDraft{}
	contentDraft := &model.CourseContentDraft{
		CourseIntroduce: req.Introduce,
		CourseDetail:    req.Detail,
		UsePeople:       req.UsePeople,
	}
	draft.CoverUrl = req.CoverUrl
	draft.PurchaseEndTime = parseTimeOrNow(req.PurchaseEndTime)

	if course == nil {
		draft.Price = int64(req.Price)
		draft.ValidDuration = int64(req.ValidDuration)
		draft.Status = CourseStatusNoUpShelf
		if len(categoryIdList) == 3 {
			draft.FirstCateId = categoryIdList[0]
			draft.SecondCateId = categoryIdList[1]
			draft.ThirdCateId = categoryIdList[2]
		}
		draft.Free = boolToInt(req.Free)
		draft.Name = req.Name
	}

	if req.Id == 0 {
		id := idgen.NextID()
		draft.Id = id
		draft.Step = CourseStepBaseInfo
		draft.Creater = userId
		draft.Updater = userId
		contentDraft.Id = id
		contentDraft.Creater = userId
		contentDraft.Updater = userId
		if err := s.courseDraftModel.Insert(ctx, draft); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "新增课程草稿失败")
		}
		if err := s.courseContentDraftModel.Insert(ctx, contentDraft); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "新增课程内容草稿失败")
		}
		return &types.CourseSaveVO{Id: id}, nil
	}

	draft.Id = req.Id
	draft.Updater = userId
	if err := s.courseDraftModel.UpdateById(ctx, draft,
		"name", "cover_url", "price", "valid_duration", "free", "status",
		"first_cate_id", "second_cate_id", "third_cate_id", "purchase_end_time"); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新课程草稿失败")
	}
	contentDraft.Id = req.Id
	contentDraft.Updater = userId
	if err := s.courseContentDraftModel.UpdateById(ctx, contentDraft); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新课程内容草稿失败")
	}
	return &types.CourseSaveVO{Id: req.Id}, nil
}

// checkSameName 校验同名课程（正式 + 草稿，排除自身）。
func (s *courseService) checkSameName(ctx context.Context, id int64, name string) error {
	num, err := s.courseModel.CountSameName(ctx, name)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "校验课程名称失败")
	}
	if num > 0 {
		return xerr.BadRequestf("课程名称已存在")
	}
	num, err = s.courseDraftModel.CountByName(ctx, name, id)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "校验课程名称失败")
	}
	if num > 0 {
		return xerr.BadRequestf("课程名称已存在")
	}
	return nil
}

// checkCategory 校验三级分类并返回一/二/三级分类 id 列表。
func (s *courseService) checkCategory(ctx context.Context, thirdCateId int64) ([]int64, error) {
	third, err := s.categoryModel.FindById(ctx, thirdCateId)
	if err == sql.ErrNoRows {
		return nil, xerr.BadRequestf("课程分类不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	if third.Level != 3 {
		return nil, xerr.BadRequestf("课程必须选择三级分类")
	}
	second, err := s.categoryModel.FindById(ctx, third.ParentId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	first, err := s.categoryModel.FindById(ctx, second.ParentId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	return []int64{first.Id, second.Id, third.Id}, nil
}

// GetBaseInfo 获取课程基础信息。
func (s *courseService) GetBaseInfo(ctx context.Context, id int64, see bool) (*types.CourseBaseInfoVO, error) {
	var vo *types.CourseBaseInfoVO
	if see {
		course, err := s.courseModel.FindById(ctx, id)
		if err != nil && err != sql.ErrNoRows {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
		}
		if course != nil {
			content, err := s.courseContentModel.FindById(ctx, id)
			if err != nil && err != sql.ErrNoRows {
				return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程内容失败")
			}
			vo = &types.CourseBaseInfoVO{
				Id:            course.Id,
				Name:          course.Name,
				CoverUrl:      course.CoverUrl,
				Price:         int32(course.Price),
				ValidDuration: int32(course.ValidDuration),
				Free:          intToBool(course.Free),
				FirstCateId:   course.FirstCateId,
				SecondCateId:  course.SecondCateId,
				ThirdCateId:   course.ThirdCateId,
				Status:        int32(course.Status),
				Step:          int32(course.Step),
				CanUpdate:     true,
				CataTotalNum:  int32(course.SectionNum.Int64),
				CoureScore:    float64(course.Score.Int64) / 10,
				CreateTime:    course.CreateTime.Format(timeFormat),
				UpdateTime:    course.UpdateTime.Format(timeFormat),
				Creater:       course.Creater,
				Updater:       course.Updater,
				Score:         int32(course.Score.Int64),
				PurchaseStartTime: formatNullTime(course.PurchaseStartTime),
				PurchaseEndTime:   course.PurchaseEndTime.Format(timeFormat),
			}
			if content != nil {
				vo.Detail = content.CourseDetail
				vo.Introduce = content.CourseIntroduce
				vo.UsePeople = content.UsePeople
			}
		}
	}
	if vo == nil {
		draft, err := s.courseDraftModel.FindById(ctx, id)
		if err == sql.ErrNoRows {
			return &types.CourseBaseInfoVO{}, nil
		}
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
		}
		contentDraft, err := s.courseContentDraftModel.FindById(ctx, id)
		if err != nil && err != sql.ErrNoRows {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程内容草稿失败")
		}
		vo = &types.CourseBaseInfoVO{
			Id:            draft.Id,
			Name:          draft.Name,
			CoverUrl:      draft.CoverUrl,
			Price:         int32(draft.Price),
			ValidDuration: int32(draft.ValidDuration),
			Free:          intToBool(draft.Free),
			FirstCateId:   draft.FirstCateId,
			SecondCateId:  draft.SecondCateId,
			ThirdCateId:   draft.ThirdCateId,
			Status:        int32(draft.Status),
			Step:          int32(draft.Step),
			CanUpdate:     intToBool(draft.CanUpdate),
			CataTotalNum:  int32(draft.SectionNum),
			CreateTime:    draft.CreateTime.Format(timeFormat),
			UpdateTime:    draft.UpdateTime.Format(timeFormat),
			Creater:       draft.Creater,
			Updater:       draft.Updater,
			PurchaseStartTime: formatNullTime(draft.PurchaseStartTime),
			PurchaseEndTime:   draft.PurchaseEndTime.Format(timeFormat),
		}
		if contentDraft != nil {
			vo.Detail = contentDraft.CourseDetail
			vo.Introduce = contentDraft.CourseIntroduce
			vo.UsePeople = contentDraft.UsePeople
		}
	}

	// 创建者/更新者姓名
	operatorIds := make([]int64, 0, 2)
	if vo.Creater > 0 {
		operatorIds = append(operatorIds, vo.Creater)
	}
	if vo.Updater > 0 && vo.Updater != vo.Creater {
		operatorIds = append(operatorIds, vo.Updater)
	}
	if len(operatorIds) > 0 {
		userMap, err := s.userDetailModel.FindByIds(ctx, operatorIds)
		if err == nil {
			if u, ok := userMap[vo.Creater]; ok {
				vo.CreaterName = u.Name
			}
			if u, ok := userMap[vo.Updater]; ok {
				vo.UpdaterName = u.Name
			}
		}
	}
	// 分类名称
	cateNames, err := s.categoryNames(ctx, vo.FirstCateId, vo.SecondCateId, vo.ThirdCateId)
	if err == nil {
		vo.CateNames = cateNames
	}
	return vo, nil
}

// categoryNames 拼接一/二/三级分类名称。
func (s *courseService) categoryNames(ctx context.Context, first, second, third int64) (string, error) {
	ids := make([]int64, 0, 3)
	for _, id := range []int64{first, second, third} {
		if id > 0 {
			ids = append(ids, id)
		}
	}
	if len(ids) == 0 {
		return "", nil
	}
	cates, err := s.categoryModel.FindByIds(ctx, ids)
	if err != nil {
		return "", err
	}
	nameMap := make(map[int64]string, len(cates))
	for _, c := range cates {
		nameMap[c.Id] = c.Name
	}
	return strings.Join([]string{nameMap[first], nameMap[second], nameMap[third]}, "/"), nil
}

// CheckName 校验课程名称是否已存在。
func (s *courseService) CheckName(ctx context.Context, name string, id int64) (*types.NameExistVO, error) {
	vo := &types.NameExistVO{Existed: false}
	num, err := s.courseModel.CountSameName(ctx, name)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "校验课程名称失败")
	}
	if num > 0 {
		vo.Existed = true
		return vo, nil
	}
	num, err = s.courseDraftModel.CountByName(ctx, name, id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "校验课程名称失败")
	}
	vo.Existed = num > 0
	return vo, nil
}

// Delete 删除课程及关联数据。
func (s *courseService) Delete(ctx context.Context, id int64) error {
	if err := s.courseDraftModel.DeleteById(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程草稿失败")
	}
	if err := s.courseContentDraftModel.DeleteById(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程内容草稿失败")
	}
	if err := s.cataSubjectDraftModel.DeleteByCourseId(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程题目关系失败")
	}
	if err := s.cataDraftModel.DeleteByCourseId(ctx, id, []int64{CataTypeChapter, CataTypeSection, CataTypePractice}); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程目录草稿失败")
	}
	if err := s.teacherDraftModel.DeleteByCourseId(ctx, id); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除课程老师关系失败")
	}
	return nil
}

// QueryForPage 管理端课程分页搜索。
func (s *courseService) QueryForPage(ctx context.Context, req *types.CoursePageQuery) (*result.Page, error) {
	offset, limit := req.Normalize()
	if req.Status == CourseStatusNoUpShelf || req.Status == CourseStatusDownShelf {
		return s.draftPage(ctx, req, offset, limit)
	}
	return s.formalPage(ctx, req, offset, limit)
}

func (s *courseService) draftPage(ctx context.Context, req *types.CoursePageQuery, offset, limit int64) (*result.Page, error) {
	cond := &model.CourseDraftPageCond{
		Keyword:     req.Keyword,
		FirstCateId: req.FirstCateId,
		SecondCateId: req.SecondCateId,
		ThirdCateId: req.ThirdCateId,
		CourseType:  req.CourseType,
		Status:      req.Status,
		BeginTime:   req.BeginTime,
		EndTime:     req.EndTime,
		OrderBy:     sortBy(req.SortBy, req.IsAsc),
		Offset:      offset,
		Limit:       limit,
	}
	list, total, err := s.courseDraftModel.Page(ctx, cond)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程草稿失败")
	}
	if len(list) == 0 {
		return &result.Page{List: []*types.CoursePageVO{}, Total: 0, Pages: 0}, nil
	}
	vos, err := s.fillPageVO(ctx, list)
	if err != nil {
		return nil, err
	}
	return &result.Page{List: vos, Total: total, Pages: page.CalcPages(total, limit)}, nil
}

func (s *courseService) formalPage(ctx context.Context, req *types.CoursePageQuery, offset, limit int64) (*result.Page, error) {
	cond := &model.CoursePageCond{
		Keyword:     req.Keyword,
		FirstCateId: req.FirstCateId,
		SecondCateId: req.SecondCateId,
		ThirdCateId: req.ThirdCateId,
		CourseType:  req.CourseType,
		Status:      req.Status,
		BeginTime:   req.BeginTime,
		EndTime:     req.EndTime,
		OrderBy:     sortBy(req.SortBy, req.IsAsc),
		Offset:      offset,
		Limit:       limit,
	}
	list, total, err := s.courseModel.Page(ctx, cond)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	if len(list) == 0 {
		return &result.Page{List: []*types.CoursePageVO{}, Total: 0, Pages: 0}, nil
	}
	vos, err := s.fillFormalPageVO(ctx, list)
	if err != nil {
		return nil, err
	}
	return &result.Page{List: vos, Total: total, Pages: page.CalcPages(total, limit)}, nil
}

func sortBy(sortBy string, isAsc bool) string {
	if sortBy == "" {
		return ""
	}
	order := "ASC"
	if !isAsc {
		order = "DESC"
	}
	return sortBy + " " + order
}

// fillPageVO 组装草稿分页 VO。
func (s *courseService) fillPageVO(ctx context.Context, list []*model.CourseDraft) ([]*types.CoursePageVO, error) {
	updaterIds := make([]int64, 0, len(list))
	cateIds := make([]int64, 0, len(list)*3)
	for _, c := range list {
		updaterIds = append(updaterIds, c.Updater)
		cateIds = append(cateIds, c.FirstCateId, c.SecondCateId, c.ThirdCateId)
	}
	userMap, err := s.userDetailModel.FindByIds(ctx, distinctIDs(updaterIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户信息失败")
	}
	cateMap, err := s.categoryNameMap(ctx, distinctIDs(cateIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	vos := make([]*types.CoursePageVO, 0, len(list))
	for _, c := range list {
		v := &types.CoursePageVO{
			Id:            c.Id,
			Name:          c.Name,
			CoverUrl:      c.CoverUrl,
			Categories:    formatCateNames(cateMap, c.FirstCateId, c.SecondCateId, c.ThirdCateId),
			Price:         c.Price,
			Status:        int32(c.Status),
			Step:          int32(c.Step),
			Score:         int32(c.Score.Int64),
			Sections:      int32(c.SectionNum),
			UpdateTime:    c.UpdateTime.Format(timeFormat),
		}
		if u, ok := userMap[c.Updater]; ok {
			v.UpdaterName = u.Name
		}
		if c.PublishTime.Valid {
			v.PublishTime = c.PublishTime.Time.Format(timeFormat)
		}
		v.PurchaseEndTime = c.PurchaseEndTime.Format(timeFormat)
		vos = append(vos, v)
	}
	return vos, nil
}

func (s *courseService) fillFormalPageVO(ctx context.Context, list []*model.Course) ([]*types.CoursePageVO, error) {
	updaterIds := make([]int64, 0, len(list))
	cateIds := make([]int64, 0, len(list)*3)
	for _, c := range list {
		updaterIds = append(updaterIds, c.Updater)
		cateIds = append(cateIds, c.FirstCateId, c.SecondCateId, c.ThirdCateId)
	}
	userMap, err := s.userDetailModel.FindByIds(ctx, distinctIDs(updaterIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户信息失败")
	}
	cateMap, err := s.categoryNameMap(ctx, distinctIDs(cateIds))
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程分类失败")
	}
	vos := make([]*types.CoursePageVO, 0, len(list))
	for _, c := range list {
		v := &types.CoursePageVO{
			Id:            c.Id,
			Name:          c.Name,
			CoverUrl:      c.CoverUrl,
			Categories:    formatCateNames(cateMap, c.FirstCateId, c.SecondCateId, c.ThirdCateId),
			Price:         c.Price,
			Status:        int32(c.Status),
			Step:          int32(c.Step),
			Score:         int32(c.Score.Int64),
			Sections:      int32(c.SectionNum.Int64),
			UpdateTime:    c.UpdateTime.Format(timeFormat),
		}
		if u, ok := userMap[c.Updater]; ok {
			v.UpdaterName = u.Name
		}
		if c.PublishTime.Valid {
			v.PublishTime = c.PublishTime.Time.Format(timeFormat)
		}
		v.PurchaseEndTime = c.PurchaseEndTime.Format(timeFormat)
		vos = append(vos, v)
	}
	return vos, nil
}

func (s *courseService) categoryNameMap(ctx context.Context, ids []int64) (map[int64]string, error) {
	cates, err := s.categoryModel.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	m := make(map[int64]string, len(cates))
	for _, c := range cates {
		m[c.Id] = c.Name
	}
	return m, nil
}

func formatCateNames(m map[int64]string, first, second, third int64) string {
	return strings.Join([]string{m[first], m[second], m[third]}, "/")
}

// QueryCourseIdByName 按名称模糊查询课程 id 列表。
func (s *courseService) QueryCourseIdByName(ctx context.Context, name string) ([]int64, error) {
	list, err := s.courseModel.FindByName(ctx, name)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	ids := make([]int64, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.Id)
	}
	sort.Slice(ids, func(i, j int) bool { return ids[i] < ids[j] })
	return ids, nil
}

// QuerySimpleInfoList 按条件列表查询课程简要信息。
func (s *courseService) QuerySimpleInfoList(ctx context.Context, ids, thirdCataIds []int64) ([]*types.CourseSimpleInfoDTO, error) {
	if len(ids) == 0 && len(thirdCataIds) == 0 {
		return []*types.CourseSimpleInfoDTO{}, nil
	}
	var list []*model.Course
	var err error
	if len(ids) > 0 {
		list, err = s.courseModel.FindByIds(ctx, ids)
	} else {
		all, err2 := s.courseModel.ListAll(ctx)
		if err2 != nil {
			return nil, xerr.Wrap(err2, xerr.CodeInternal, "查询课程失败")
		}
		set := make(map[int64]bool, len(thirdCataIds))
		for _, id := range thirdCataIds {
			set[id] = true
		}
		for _, c := range all {
			if set[c.ThirdCateId] {
				list = append(list, c)
			}
		}
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程失败")
	}
	result := make([]*types.CourseSimpleInfoDTO, 0, len(list))
	for _, c := range list {
		result = append(result, &types.CourseSimpleInfoDTO{
			Id:              c.Id,
			Name:            c.Name,
			CoverUrl:        c.CoverUrl,
			Price:           int32(c.Price),
			Free:            intToBool(c.Free),
			Status:          int32(c.Status),
			ValidDuration:   int32(c.ValidDuration),
			SectionNum:      int32(c.SectionNum.Int64),
			PurchaseEndTime: c.PurchaseEndTime.Format(timeFormat),
			FirstCateId:     c.FirstCateId,
			SecondCateId:    c.SecondCateId,
			ThirdCateId:     c.ThirdCateId,
		})
	}
	return result, nil
}
