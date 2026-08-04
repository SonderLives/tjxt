package logic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/response"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 支付渠道列表 ============

type PayChannelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelsLogic {
	return &PayChannelsLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayChannelsLogic) PayChannels() (resp *result.R, err error) {
	channels, err := l.svcCtx.PayService.ListChannels(l.ctx)
	if err != nil {
		return nil, err
	}
	return result.OkData(channels), nil
}

// ============ 支付申请 ============

type PayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PayOrderLogic) PayOrder(req *types.PayApplyReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	payURL, err := l.svcCtx.PayService.ApplyPay(l.ctx, userID, req.OrderId, req.PayChannelCode)
	if err != nil {
		return nil, err
	}
	return result.OkData(payURL), nil
}
