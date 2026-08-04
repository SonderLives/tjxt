package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPayChannelByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryPayChannelByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPayChannelByCodeLogic {
	return &QueryPayChannelByCodeLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QueryPayChannelByCodeLogic) QueryPayChannelByCode(in *pb.QueryPayChannelByCodeRequest) (*pb.PayChannelResponse, error) {
	if in.ChannelCode == "" {
		return nil, xerr.BadRequestf("渠道编码不能为空")
	}
	m, err := l.svcCtx.PayChannelModel.FindByCode(l.ctx, in.ChannelCode)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("支付渠道不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付渠道失败")
	}
	return toPayChannelResp(m), nil
}