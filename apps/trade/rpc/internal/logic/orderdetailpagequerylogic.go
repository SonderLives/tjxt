package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailPageQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailPageQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPageQueryLogic {
	return &OrderDetailPageQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailPageQueryLogic) OrderDetailPageQuery(in *pb.OrderDetailPageRequest) (*pb.OrderDetailPageReply, error) {
	// todo: add your logic here and delete this line

	return &pb.OrderDetailPageReply{}, nil
}
