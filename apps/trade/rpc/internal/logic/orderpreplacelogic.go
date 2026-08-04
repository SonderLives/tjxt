package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPrePlaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderPrePlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPrePlaceLogic {
	return &OrderPrePlaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 订单 =====
func (l *OrderPrePlaceLogic) OrderPrePlace(in *pb.PrePlaceRequest) (*pb.OrderConfirmVO, error) {
	// todo: add your logic here and delete this line

	return &pb.OrderConfirmVO{}, nil
}
