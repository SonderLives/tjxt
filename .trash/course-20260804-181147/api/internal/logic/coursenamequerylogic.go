// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/pkg/response"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseNameQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseNameQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseNameQueryLogic {
	return &CourseNameQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseNameQueryLogic) CourseNameQuery(req *types.NameQuery) (resp *result.R, err error) {
	ids, err := l.svcCtx.CourseService.QueryCourseIdByName(l.ctx, req.Name)
	if err != nil {
		return nil, err
	}
	return result.OkData(ids), nil
}
