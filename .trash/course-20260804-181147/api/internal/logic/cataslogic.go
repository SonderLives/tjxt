package logic

import (
	"context"

	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/response"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 保存章节 ============

type CourseCatasSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCatasSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatasSaveLogic {
	return &CourseCatasSaveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseCatasSaveLogic) SaveCatas(id, step int64, list []*types.CataSaveDTO) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.SaveCatas(l.ctx, id, step, list); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 获取课程的章节 ============

type CourseCatasLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCatasLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatasLogic {
	return &CourseCatasLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseCatasLogic) Catas(id int64, req *types.CourseCatasQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.QueryCatas(l.ctx, id, req.See, req.WithPractice)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 课程视频（保存小节媒资） ============

type CourseMediaSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseMediaSaveLogic {
	return &CourseMediaSaveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseMediaSaveLogic) SaveMedia(id int64, list []*types.CourseMediaSaveDTO) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.SaveMedia(l.ctx, id, list); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 保存小节或练习中的题目 ============

type CourseSubjectsSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsSaveLogic {
	return &CourseSubjectsSaveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSubjectsSaveLogic) SaveSubjects(id int64, list []*types.CataSubjectDTO) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.SaveSubjects(l.ctx, id, list); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 获取小节或练习中的题目 ============

type CourseSubjectsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSubjectsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSubjectsLogic {
	return &CourseSubjectsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSubjectsLogic) Subjects(id int64) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.GetSubjects(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 根据课程id，查询所有章节的序号 ============

type CourseCataIndexLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataIndexLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataIndexLogic {
	return &CourseCataIndexLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseCataIndexLogic) CataIndexList(courseId int64) (resp *result.R, err error) {
	list, err := l.svcCtx.CatalogueService.GetCatasIndexList(l.ctx, courseId)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 生成练习id ============

type CourseGeneratorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseGeneratorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseGeneratorLogic {
	return &CourseGeneratorLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseGeneratorLogic) Generator() (resp *result.R, err error) {
	id := idgen.NextID()
	return result.OkData(&types.CourseCataIdVO{Id: id}), nil
}
