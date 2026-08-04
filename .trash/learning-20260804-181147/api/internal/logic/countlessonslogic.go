// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"fmt"

	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CountLessonsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCountLessonsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CountLessonsLogic {
	return &CountLessonsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CountLessonsLogic) CountLessons(req *types.LessonRequest) (resp *types.Result, err error) {
	if req.CourseId <= 0 {
		return nil, fmt.Errorf("courseId must be positive")
	}
	count, err := l.svcCtx.LessonService.CountLessons(l.ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	return success(count), nil
}
