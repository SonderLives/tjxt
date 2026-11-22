package paylogic

import (
	"context"
	"errors"

	"tjxt/common/errorx"
	"tjxt/service/pay/rpc/internal/model"
	"tjxt/service/pay/rpc/internal/svc"
	"tjxt/service/pay/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePayChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePayChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePayChannelLogic {
	return &UpdatePayChannelLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *UpdatePayChannelLogic) UpdatePayChannel(in *pb.PayChannelRequest) (*pb.EmptyResponse, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("渠道 id 非法")
	}
	exist, err := l.svcCtx.PayChannelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("支付渠道不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付渠道失败")
	}

	// 不允许修改 code，以避免与已有订单对不上
	if in.ChannelCode != "" && in.ChannelCode != exist.ChannelCode {
		return nil, errorx.BadRequestf("渠道编码不允许修改")
	}

	if in.Name != "" {
		exist.Name = in.Name
	}
	if in.ChannelPriority > 0 {
		exist.ChannelPriority = int64(in.ChannelPriority)
	}
	if in.ChannelIcon != "" {
		exist.ChannelIcon = in.ChannelIcon
	}

	if err := l.svcCtx.PayChannelModel.Update(l.ctx, exist); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新支付渠道失败")
	}
	return &pb.EmptyResponse{}, nil
}
