package paylogic

import (
	"context"
	"errors"

	"tjxt/pkg/errorx"
	"tjxt/rpc/pay/internal/model"
	"tjxt/rpc/pay/internal/svc"
	"tjxt/rpc/pay/pb"

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
		return nil, errorx.BadRequestf("渠道编码不能为空")
	}
	m, err := l.svcCtx.PayChannelModel.FindByCode(l.ctx, in.ChannelCode)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("支付渠道不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付渠道失败")
	}
	return toPayChannelResp(m), nil
}
