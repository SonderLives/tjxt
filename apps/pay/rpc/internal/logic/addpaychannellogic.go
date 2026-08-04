package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("新增渠道不应携带 id")
	}
	if in.Name == "" || in.ChannelCode == "" {
		return nil, xerr.BadRequestf("渠道名称与编码不能为空")
	}
	if _, err := l.svcCtx.PayChannelModel.FindByCode(l.ctx, in.ChannelCode); err == nil {
		return nil, xerr.Conflict("渠道编码已存在")
	} else if !isNotFound(err) {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询渠道编码失败")
	}

	po := buildPayChannel(in)
	res, err := l.svcCtx.PayChannelModel.Insert(l.ctx, po)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "新增支付渠道失败")
	}
	id, _ := res.LastInsertId()
	return &pb.PayChannelIdResponse{Id: id}, nil
}