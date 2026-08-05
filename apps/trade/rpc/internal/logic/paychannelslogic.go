package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelsLogic {
	return &PayChannelsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelsLogic) PayChannels(in *pb.Empty) (*pb.PayChannelVOList, error) {
	resp, err := l.svcCtx.PayRpc.ListPayChannels(l.ctx, &payclient.ListPayChannelsRequest{})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询支付渠道失败")
	}

	items := make([]*pb.PayChannelVO, 0, len(resp.List))
	for _, ch := range resp.List {
		items = append(items, &pb.PayChannelVO{
			Id:              ch.Id,
			Name:            ch.Name,
			ChannelCode:     ch.ChannelCode,
			ChannelIcon:     ch.ChannelIcon,
			ChannelPriority: ch.ChannelPriority,
		})
	}
	return &pb.PayChannelVOList{Items: items}, nil
}
