package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyGetLogic {
	return &RefundApplyGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyGetLogic) RefundApplyGet(in *pb.IdRequest) (*pb.RefundApplyVO, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundApplyVO{}, nil
}
