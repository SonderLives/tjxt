// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"common/result"
	"course/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseDTOByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseDTOByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseDTOByIdLogic {
	return &CourseDTOByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseDTOByIdLogic) CourseDTOById(id int64) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.GetCourseDTOById(l.ctx, id)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}
