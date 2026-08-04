package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.PlaceOrderResultVO{}, nil
}
