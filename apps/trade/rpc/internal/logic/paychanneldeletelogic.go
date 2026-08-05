package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelDeleteLogic {
	return &PayChannelDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelDeleteLogic) PayChannelDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("支付渠道ID不能为空")
	}

	// 渠道无物理删除，置为禁用（status=2）等价逻辑删除
	if _, err := l.svcCtx.PayRpc.UpdatePayChannelStatus(l.ctx, &payclient.UpdatePayChannelStatusRequest{
		Id:     in.Id,
		Status: 2,
	}); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除支付渠道失败")
	}
	return &pb.Empty{}, nil
}
