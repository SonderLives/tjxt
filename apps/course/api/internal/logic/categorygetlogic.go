// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CategoryGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CategoryGetLogic {
	return &CategoryGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryGetLogic) CategoryGet(req *types.IdPathReq) (resp *types.CategoryInfoVO, err error) {
	// todo: add your logic here and delete this line

	return
}
