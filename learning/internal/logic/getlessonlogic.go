// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"learning/internal/svc"
	"learning/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLessonLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLessonLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLessonLogic {
	return &GetLessonLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLessonLogic) GetLesson(req *types.LessonRequest) (resp *types.Result, err error) {
	userID, err := currentUserID(l.ctx)
	if err != nil {
		return nil, err
	}
	lesson, err := l.svcCtx.LessonService.GetLesson(l.ctx, userID, req.CourseId)
	if err != nil {
		return nil, err
	}
	return success(lessonResponse(lesson)), nil
}
