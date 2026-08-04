// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPrePlaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPrePlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPrePlaceLogic {
	return &OrderPrePlaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderPrePlaceLogic) OrderPrePlace(req *types.PrePlaceOrderReq) (resp *types.OrderConfirmVO, err error) {
	// todo: add your logic here and delete this line

	return
}
