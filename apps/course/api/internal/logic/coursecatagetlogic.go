// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CourseCataGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCourseCataGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CourseCataGetLogic {
	return &CourseCataGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CourseCataGetLogic) CourseCataGet(req *types.CourseCataQueryReq) (resp []types.CatalogueVO, err error) {
	// todo: add your logic here and delete this line

	return
}
