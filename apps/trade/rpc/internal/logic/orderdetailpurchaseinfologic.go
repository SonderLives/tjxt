package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailPurchaseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailPurchaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPurchaseInfoLogic {
	return &OrderDetailPurchaseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailPurchaseInfoLogic) OrderDetailPurchaseInfo(in *pb.PurchaseInfoRequest) (*pb.PurchaseInfoReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PurchaseInfoReply{}, nil
}
