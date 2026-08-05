package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("渠道名称不能为空")
	}
	if in.ChannelCode == "" {
		return nil, xerr.BadRequestf("渠道编码不能为空")
	}

	// trade 无本地 pay_channel 表，统一委托 pay 服务
	resp, err := l.svcCtx.PayRpc.AddPayChannel(l.ctx, &payclient.PayChannelRequest{
		Name:            in.Name,
		ChannelCode:     in.ChannelCode,
		ChannelIcon:     in.ChannelIcon,
		ChannelPriority: in.ChannelPriority,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增支付渠道失败")
	}
	return &pb.IdReply{Id: resp.Id}, nil
}
