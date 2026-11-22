// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/api/trade/backend/logic"
	"tjxt/api/trade/backend/pb"
	"tjxt/api/trade/internal/svc"
	"tjxt/api/trade/internal/types"

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
	reply, err := backendlogic.NewOrderPlaceLogic(l.ctx, l.svcCtx.Backend).OrderPlace(&pb.PlaceOrderRequest{
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
