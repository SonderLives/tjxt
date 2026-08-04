package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayResultQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayResultQueryLogic {
	return &PayResultQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayResultQueryLogic) PayResultQuery(in *pb.PayResultQueryRequest) (*pb.PayResultDTO, error) {
	// todo: add your logic here and delete this line

	return &pb.PayResultDTO{}, nil
}
