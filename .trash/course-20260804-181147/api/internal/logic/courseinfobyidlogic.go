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

type CourseInfoByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseInfoByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseInfoByIdLogic {
	return &CourseInfoByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseInfoByIdLogic) CourseInfoById(id int64, req *types.CourseInfoQuery) (resp *result.R, err error) {
	vo, err := l.svcCtx.CourseService.GetInfoById(l.ctx, id, req.WithCatalogue, req.WithTeachers)
	if err != nil {
		return nil, err
	}
	return result.OkData(vo), nil
}
