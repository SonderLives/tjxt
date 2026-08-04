// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartAddLogic) CartAdd(req *types.CartAddReq) (resp *types.NamePlaceVO, err error) {
	// todo: add your logic here and delete this line

	return
}
