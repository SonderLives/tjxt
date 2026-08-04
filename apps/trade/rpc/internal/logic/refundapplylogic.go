package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyLogic {
	return &RefundApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 退款 =====
func (l *RefundApplyLogic) RefundApply(in *pb.RefundApplyRequest) (*pb.RefundResultDTO, error) {
	// todo: add your logic here and delete this line

	return &pb.RefundResultDTO{}, nil
}
