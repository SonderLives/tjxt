// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CatalogueSectionInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCatalogueSectionInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CatalogueSectionInfoLogic {
	return &CatalogueSectionInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CatalogueSectionInfoLogic) CatalogueSectionInfo(req *types.IdPathReq) (resp *types.CourseSectionInfoVO, err error) {
	// todo: add your logic here and delete this line

	return
}
