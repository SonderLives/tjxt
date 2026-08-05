package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderCancelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderCancelLogic {
	return &OrderCancelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderCancelLogic) OrderCancel(in *pb.IdRequest) (*pb.Empty, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}

	if err := l.svcCtx.OrderModel.UpdateStatus(l.ctx, in.Id, OrderStatusClosed, "用户取消"); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "取消订单失败")
	}

	details, err := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}
	for _, d := range details {
		if e := l.svcCtx.OrderDetailModel.UpdateStatus(l.ctx, d.Id, DetailStatusClosed); e != nil {
			return nil, xerr.Wrap(e, xerr.CodeInternal, "关闭订单明细失败")
		}
	}

	return &pb.Empty{}, nil
}
