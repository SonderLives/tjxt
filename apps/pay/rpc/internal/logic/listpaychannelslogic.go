package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPayChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPayChannelsLogic {
	return &ListPayChannelsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ListPayChannelsLogic) ListPayChannels(in *pb.ListPayChannelsRequest) (*pb.ListPayChannelsResponse, error) {
	list, err := l.svcCtx.PayChannelModel.FindAllEnabled(l.ctx)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付渠道失败")
	}
	resp := &pb.ListPayChannelsResponse{List: make([]*pb.PayChannelResponse, 0, len(list))}
	for _, item := range list {
		resp.List = append(resp.List, toPayChannelResp(item))
	}
	return resp, nil
}