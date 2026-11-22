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

type AddPayChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddPayChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddPayChannelLogic {
	return &AddPayChannelLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *AddPayChannelLogic) AddPayChannel(in *pb.PayChannelRequest) (*pb.PayChannelIdResponse, error) {
	if in.Id != 0 {
		return nil, errorx.BadRequestf("新增渠道不应携带 id")
	}
	if in.Name == "" || in.ChannelCode == "" {
		return nil, errorx.BadRequestf("渠道名称与编码不能为空")
	}
	if _, err := l.svcCtx.PayChannelModel.FindByCode(l.ctx, in.ChannelCode); err == nil {
		return nil, errorx.Conflict("渠道编码已存在")
	} else if !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询渠道编码失败")
	}

	po := buildPayChannel(in)
	res, err := l.svcCtx.PayChannelModel.Insert(l.ctx, po)
	if err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "新增支付渠道失败")
	}
	id, _ := res.LastInsertId()
	return &pb.PayChannelIdResponse{Id: id}, nil
}
