package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyDetailLogic {
	return &RefundApplyDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyDetailLogic) RefundApplyDetail(in *pb.IdRequest) (*pb.RefundApplyVO, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundApplyVO{}, nil
}
