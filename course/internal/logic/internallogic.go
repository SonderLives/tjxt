package logic

import (
	"context"

	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 内部调用-获取课程信息 ============

type CourseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoLogic {
	return &CourseInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseInfoLogic) Info(id int64, req *types.CourseInfoQuery) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.GetInfoById(l.ctx, id, req.WithCatalogue, req.WithTeachers)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 内部调用-课程上架入索引库信息 ============

type CourseSearchInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSearchInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSearchInfoLogic {
	return &CourseSearchInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSearchInfoLogic) SearchInfo(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.GetCourseDTOById(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 内部调用-通过老师id获取课程与出题数量 ============

type CourseInfoByTeacherIdsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseInfoByTeacherIdsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByTeacherIdsLogic {
	return &CourseInfoByTeacherIdsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseInfoByTeacherIdsLogic) InfoByTeacherIds(req *types.TeacherIdsQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.CountSubjectNumAndCourseNumOfTeacher(l.ctx, req.TeacherIds)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 内部调用-媒资被引用情况 ============

type CourseMediaUseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseMediaUseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaUseInfoLogic {
	return &CourseMediaUseInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseMediaUseInfoLogic) MediaUseInfo(req *types.MediaIdsQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CatalogueService.CountMediaUserInfo(l.ctx, req.MediaIds)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 内部调用-按名称查询课程id列表 ============

type CourseIdByNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseIdByNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseIdByNameLogic {
	return &CourseIdByNameLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseIdByNameLogic) IdByName(req *types.NameQuery) (resp *result.R, err error) {
	ids, err := l.svcCtx.CourseService.QueryCourseIdByName(l.ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return result.OkData(ids), nil
}

// ============ 内部调用-小节信息 ============

type CourseSectionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSectionInfoLogic {
	return &CourseSectionInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSectionInfoLogic) SectionInfo(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CatalogueService.GetSimpleSectionInfo(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 查询课程基本信息、目录、学习进度 ============

type CourseAndCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseAndCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseAndCatalogLogic {
	return &CourseAndCatalogLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseAndCatalogLogic) CourseAndCatalog(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.QueryCourseAndCatalog(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}
