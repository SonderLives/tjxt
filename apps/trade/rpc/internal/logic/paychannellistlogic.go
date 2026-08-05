package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelListLogic {
	return &PayChannelListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayChannelListLogic) PayChannelList(in *pb.Empty) (*pb.PayChannelListReply, error) {
	resp, err := l.svcCtx.PayRpc.ListPayChannels(l.ctx, &payclient.ListPayChannelsRequest{})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询支付渠道列表失败")
	}

	items := make([]*pb.PayChannelDTO, 0, len(resp.List))
	for _, ch := range resp.List {
		items = append(items, &pb.PayChannelDTO{
			Id:              ch.Id,
			Name:            ch.Name,
			ChannelCode:     ch.ChannelCode,
			ChannelIcon:     ch.ChannelIcon,
			ChannelPriority: ch.ChannelPriority,
			Status:          ch.Status,
		})
	}
	return &pb.PayChannelListReply{Items: items}, nil
}
