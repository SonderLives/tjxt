// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CurrentLessonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCurrentLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CurrentLessonLogic {
	return &CurrentLessonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CurrentLessonLogic) CurrentLesson() (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	lessons, _, err := l.svcCtx.LessonService.ListLessons(l.ctx, userID, 1, 1)
	if err != nil {
		return nil, err
	}
	if len(lessons) == 0 {
		return success(nil), nil
	}
	return success(lessonResponse(&lessons[0])), nil
}
