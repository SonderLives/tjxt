package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelGetLogic {
	return &PayChannelGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelGetLogic) PayChannelGet(in *pb.IdRequest) (*pb.PayChannelDTO, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("支付渠道ID不能为空")
	}

	// pay 服务未提供按 id 查询，先取全量再筛选
	resp, err := l.svcCtx.PayRpc.ListPayChannels(l.ctx, &payclient.ListPayChannelsRequest{})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询支付渠道失败")
	}

	for _, ch := range resp.List {
		if ch.Id != in.Id {
			continue
		}
		return &pb.PayChannelDTO{
			Id:              ch.Id,
			Name:            ch.Name,
			ChannelCode:     ch.ChannelCode,
			ChannelIcon:     ch.ChannelIcon,
			ChannelPriority: ch.ChannelPriority,
			Status:          ch.Status,
		}, nil
	}
	return nil, xerr.NotFound("支付渠道不存在")
}
