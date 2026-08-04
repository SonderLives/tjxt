// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLessonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLessonLogic {
	return &DeleteLessonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteLessonLogic) DeleteLesson(req *types.LessonRequest) (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	if err = l.svcCtx.LessonService.DeleteCourseFromLesson(l.ctx, userID, req.CourseId); err != nil {
		return nil, err
	}
	return success(nil), nil
}
