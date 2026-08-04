package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 购物车 =====
func (l *CartAddLogic) CartAdd(in *pb.CartAddRequest) (*pb.Empty, error) {
	// todo: add your logic here and delete this line

	return &pb.Empty{}, nil
}
