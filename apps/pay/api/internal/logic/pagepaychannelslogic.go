// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
	payclient "tjxt/apps/pay/rpc/pay"

	"github.com/zeromicro/go-zero/core/logx"
)

type PagePayChannelsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPagePayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PagePayChannelsLogic {
	return &PagePayChannelsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PagePayChannelsLogic) PagePayChannels(req *types.PayChannelPageReq) (resp *types.Result, err error) {
	r, err := l.svcCtx.PayRpc.PageQueryPayChannels(l.ctx, &payclient.PageQueryPayChannelsRequest{
		PageNo:      req.PageNo,
		PageSize:    req.PageSize,
		Name:        req.Name,
		ChannelCode: req.ChannelCode,
		Status:      int32(req.Status),
	})
	if err != nil {
		return nil, err
	}
	return &types.Result{Code: 200, Msg: "OK", Data: r}, nil
}