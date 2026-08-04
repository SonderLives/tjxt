package logic

import (
	"context"

	"common/auth"
	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 新增课程（保存基本信息） ============

type CourseSaveBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSaveBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSaveBaseInfoLogic {
	return &CourseSaveBaseInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSaveBaseInfoLogic) SaveBaseInfo(req *types.CourseBaseInfoSaveDTO) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	vo, err := l.svcCtx.CourseService.SaveBaseInfo(l.ctx, req, userID)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 获取课程基础信息 ============

type CourseBaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseBaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseBaseInfoLogic {
	return &CourseBaseInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseBaseInfoLogic) BaseInfo(id int64, req *types.CourseBaseInfoQuery) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.GetBaseInfo(l.ctx, id, req.See)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 校验课程名称 ============

type CourseCheckNameLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCheckNameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckNameLogic {
	return &CourseCheckNameLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseCheckNameLogic) CheckName(req *types.CheckNameQuery) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.CheckName(l.ctx, req.Name, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}

// ============ 删除课程 ============

type CourseDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDeleteLogic {
	return &CourseDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseDeleteLogic) Delete(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.Delete(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 课程分页查询 ============

type CoursePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCoursePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CoursePageLogic {
	return &CoursePageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CoursePageLogic) Page(req *types.CoursePageQuery) (resp *result.R, err error) {
	page, err := l.svcCtx.CourseService.QueryForPage(l.ctx, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 课程简要信息列表 ============

type CourseSimpleInfoListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSimpleInfoListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSimpleInfoListLogic {
	return &CourseSimpleInfoListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseSimpleInfoListLogic) SimpleInfoList(req *types.SimpleInfoListQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.QuerySimpleInfoList(l.ctx, req.Ids, req.ThirdCataIds)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 课程上架 ============

type CourseUpShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpShelfLogic {
	return &CourseUpShelfLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseUpShelfLogic) UpShelf(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.UpShelf(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 课程下架 ============

type CourseDownShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDownShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDownShelfLogic {
	return &CourseDownShelfLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseDownShelfLogic) DownShelf(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.DownShelf(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 课程上架校验 ============

type CourseCheckBeforeUpShelfLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCheckBeforeUpShelfLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCheckBeforeUpShelfLogic {
	return &CourseCheckBeforeUpShelfLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseCheckBeforeUpShelfLogic) CheckBeforeUpShelf(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.CheckBeforeUpShelf(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}
