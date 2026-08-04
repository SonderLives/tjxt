// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartGetLogic {
	return &CartGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartGetLogic) CartGet(req *types.CartIdReq) (resp *types.CartVO, err error) {
	// todo: add your logic here and delete this line

	return
}
