// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/pkg/response"
	"tjxt/apps/course/api/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatasIndexListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCatasIndexListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatasIndexListLogic {
	return &CourseCatasIndexListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseCatasIndexListLogic) CourseCatasIndexList(id int64) (resp *result.R, err error) {
	list, err := l.svcCtx.CourseService.QueryCatas(l.ctx, id, true, false)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}
