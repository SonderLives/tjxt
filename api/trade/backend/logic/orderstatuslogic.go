package tradelogic

import (
	"context"
	"errors"

	"tjxt/api/trade/backend/model"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatusLogic {
	return &OrderStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderStatusLogic) OrderStatus(in *pb.IdRequest) (*pb.PlaceOrderResultVO, error) {
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单失败")
	}

	return &pb.PlaceOrderResultVO{
		OrderId:   order.Id,
		PayAmount: order.RealAmount,
		Status:    int32(order.Status),
	}, nil
}
