package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayApplyLogic {
	return &PayApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 支付 =====
func (l *PayApplyLogic) PayApply(in *pb.PayApplyRequest) (*pb.PayApplyReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PayApplyReply{}, nil
}
