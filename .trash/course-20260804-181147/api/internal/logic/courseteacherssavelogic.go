// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/response"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseTeachersSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseTeachersSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseTeachersSaveLogic {
	return &CourseTeachersSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseTeachersSaveLogic) CourseTeachersSave(req *types.CourseTeacherSaveDTO) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.CourseService.SaveTeachers(l.ctx, req, userID); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}
