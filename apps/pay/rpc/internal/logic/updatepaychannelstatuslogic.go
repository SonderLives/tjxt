package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePayChannelStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePayChannelStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePayChannelStatusLogic {
	return &UpdatePayChannelStatusLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdatePayChannelStatusLogic) UpdatePayChannelStatus(in *pb.UpdatePayChannelStatusRequest) (*pb.EmptyResponse, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("渠道 id 非法")
	}
	if in.Status != PayChannelStatusEnabled && in.Status != PayChannelStatusDisabled {
		return nil, xerr.BadRequestf("渠道状态非法")
	}
	exist, err := l.svcCtx.PayChannelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("支付渠道不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付渠道失败")
	}
	exist.Status = int64(in.Status)
	if err := l.svcCtx.PayChannelModel.Update(l.ctx, exist); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新渠道状态失败")
	}
	return &pb.EmptyResponse{}, nil
}