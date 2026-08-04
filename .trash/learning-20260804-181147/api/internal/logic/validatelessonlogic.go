// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ValidateLessonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewValidateLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ValidateLessonLogic {
	return &ValidateLessonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ValidateLessonLogic) ValidateLesson(req *types.LessonRequest) (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	lesson, err := l.svcCtx.LessonService.GetLesson(l.ctx, userID, req.CourseId)
	if err != nil || lesson.Status == 3 {
		return success(int64(0)), nil
	}
	return success(lesson.Id), nil
}
