package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

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
			return nil, xerr.NotFound("订单不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}

	return &pb.PlaceOrderResultVO{
		OrderId:   order.Id,
		PayAmount: order.RealAmount,
		Status:    int32(order.Status),
	}, nil
}
