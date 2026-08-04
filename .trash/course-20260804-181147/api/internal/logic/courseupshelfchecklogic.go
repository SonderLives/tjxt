// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/pkg/response"
	"tjxt/apps/course/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseUpShelfCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseUpShelfCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseUpShelfCheckLogic {
	return &CourseUpShelfCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseUpShelfCheckLogic) CheckBeforeUpShelf(id int64) (resp *result.R, err error) {
	if err := l.svcCtx.CourseService.CheckBeforeUpShelf(l.ctx, id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}
