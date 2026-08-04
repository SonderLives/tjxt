// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCatalogsGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCatalogsGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCatalogsGetLogic {
	return &CourseCatalogsGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseCatalogsGetLogic) CourseCatalogsGet(req *types.IdPathReq) (resp *types.CourseAndSectionVO, err error) {
	// todo: add your logic here and delete this line

	return
}
