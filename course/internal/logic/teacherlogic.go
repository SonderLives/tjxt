package logic

import (
	"context"

	"common/auth"
	"common/result"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 保存老师信息 ============

type CourseTeacherSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeacherSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeacherSaveLogic {
	return &CourseTeacherSaveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseTeacherSaveLogic) SaveTeachers(req *types.CourseTeacherSaveDTO) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.CourseService.SaveTeachers(l.ctx, req, userID); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 查询课程相关的老师信息 ============

type CourseTeacherLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeacherLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeacherLogic {
	return &CourseTeacherLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CourseTeacherLogic) Teachers(id int64, req *types.CourseTeacherQuery) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.QueryTeachers(l.ctx, id, req.See)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}
