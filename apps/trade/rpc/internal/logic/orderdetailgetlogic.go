package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailGetLogic {
	return &OrderDetailGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 订单明细 =====
func (l *OrderDetailGetLogic) OrderDetailGet(in *pb.IdRequest) (*pb.OrderDetailAdminVO, error) {
	// todo: add your logic here and delete this line

	return &pb.OrderDetailAdminVO{}, nil
}
