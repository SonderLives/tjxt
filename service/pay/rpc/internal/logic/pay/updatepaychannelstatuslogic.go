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
		return nil, errorx.BadRequestf("渠道 id 非法")
	}
	if in.Status != PayChannelStatusEnabled && in.Status != PayChannelStatusDisabled {
		return nil, errorx.BadRequestf("渠道状态非法")
	}
	exist, err := l.svcCtx.PayChannelModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("支付渠道不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付渠道失败")
	}
	exist.Status = int64(in.Status)
	if err := l.svcCtx.PayChannelModel.Update(l.ctx, exist); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "更新渠道状态失败")
	}
	return &pb.EmptyResponse{}, nil
}
