package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartGetLogic {
	return &CartGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartGetLogic) CartGet(in *pb.IdRequest) (*pb.CartVO, error) {
	// todo: add your logic here and delete this line

	return &pb.CartVO{}, nil
}
