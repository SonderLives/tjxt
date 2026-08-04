package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPlaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderPlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPlaceLogic {
	return &OrderPlaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderPlaceLogic) OrderPlace(in *pb.PlaceOrderRequest) (*pb.PlaceOrderResultVO, error) {
	// todo: add your logic here and delete this line

	return &pb.PlaceOrderResultVO{}, nil
}
