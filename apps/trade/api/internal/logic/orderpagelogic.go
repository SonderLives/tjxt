// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPageLogic {
	return &OrderPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderPageLogic) OrderPage(req *types.OrderPageReq) (resp *types.OrderPageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
