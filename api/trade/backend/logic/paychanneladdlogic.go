package tradelogic

import (
	"context"

	payclient "tjxt/rpc/pay/client/pay"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelAddLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayChannelAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelAddLogic {
	return &PayChannelAddLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 支付渠道 =====
func (l *PayChannelAddLogic) PayChannelAdd(in *pb.PayChannelDTO) (*pb.IdReply, error) {
	if in.Name == "" {
		return nil, errorx.BadRequestf("渠道名称不能为空")
	}
	if in.ChannelCode == "" {
		return nil, errorx.BadRequestf("渠道编码不能为空")
	}

	// trade 无本地 pay_channel 表，统一委托 pay 服务
	resp, err := l.svcCtx.PayRpc.AddPayChannel(l.ctx, &payclient.PayChannelRequest{
		Name:            in.Name,
		ChannelCode:     in.ChannelCode,
		ChannelIcon:     in.ChannelIcon,
		ChannelPriority: in.ChannelPriority,
	})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "新增支付渠道失败")
	}
	return &pb.IdReply{Id: resp.Id}, nil
}
