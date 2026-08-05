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

type OrderPlaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPlaceLogic {
	return &OrderPlaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderPlaceLogic) OrderPlace(req *types.PlaceOrderReq) (resp *types.PlaceOrderResultVO, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderPlace(l.ctx, &pb.PlaceOrderRequest{
		OrderId:   req.OrderId,
		CourseIds: req.CourseIds,
		CouponIds: req.CouponIds,
	})
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
