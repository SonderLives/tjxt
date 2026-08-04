package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundResultQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundResultQueryLogic {
	return &RefundResultQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundResultQueryLogic) RefundResultQuery(in *pb.RefundResultQueryRequest) (*pb.RefundResultDTO, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundResultDTO{}, nil
}
