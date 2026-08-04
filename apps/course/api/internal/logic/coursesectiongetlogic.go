// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseSectionGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseSectionGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseSectionGetLogic {
	return &CourseSectionGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseSectionGetLogic) CourseSectionGet(req *types.IdPathReq) (resp *types.CourseSectionInfoVO, err error) {
	// todo: add your logic here and delete this line

	return
}
