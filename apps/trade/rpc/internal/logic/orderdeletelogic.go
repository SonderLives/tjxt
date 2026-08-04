package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDeleteLogic {
	return &OrderDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDeleteLogic) OrderDelete(in *pb.IdRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
