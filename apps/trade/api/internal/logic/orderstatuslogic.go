// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatusLogic {
	return &OrderStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderStatusLogic) OrderStatus(req *types.OrderIdReq) (resp *types.PlaceOrderResultVO, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderStatus(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.PlaceOrderResultVO{
		OrderId:    reply.OrderId,
		PayAmount:  reply.PayAmount,
		PayOutTime: reply.PayOutTime,
		Status:     int64(reply.Status),
	}, nil
}
