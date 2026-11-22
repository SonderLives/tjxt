package tradelogic

import (
	"context"

	payclient "tjxt/rpc/pay/client/pay"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询支付渠道失败")
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
