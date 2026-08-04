package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyNextLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyNextLogic {
	return &RefundApplyNextLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyNextLogic) RefundApplyNext(in *pb.Empty) (*pb.RefundApplyVO, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundApplyVO{}, nil
}
