package tradelogic

import (
	"context"

	payclient "tjxt/service/pay/rpc/client/pay"
	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/errorx"

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
		return nil, errorx.BadRequestf("支付渠道ID不能为空")
	}

	// pay 服务未提供按 id 查询，先取全量再筛选
	resp, err := l.svcCtx.PayRpc.ListPayChannels(l.ctx, &payclient.ListPayChannelsRequest{})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询支付渠道失败")
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
	return nil, errorx.NotFound("支付渠道不存在")
}
