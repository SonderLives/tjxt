package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PageQueryPayChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPageQueryPayChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PageQueryPayChannelsLogic {
	return &PageQueryPayChannelsLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *PageQueryPayChannelsLogic) PageQueryPayChannels(in *pb.PageQueryPayChannelsRequest) (*pb.PageQueryPayChannelsResponse, error) {
	offset, limit := normalizePage(in)
	list, total, err := l.svcCtx.PayChannelModel.PageList(l.ctx, in.Name, in.ChannelCode, int64(in.Status), offset, limit)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "分页查询支付渠道失败")
	}
	resp := &pb.PageQueryPayChannelsResponse{
		Total: total,
		Pages: page.CalcPages(total, limit),
		List:  make([]*pb.PayChannelResponse, 0, len(list)),
	}
	for _, item := range list {
		resp.List = append(resp.List, toPayChannelResp(item))
	}
	return resp, nil
}