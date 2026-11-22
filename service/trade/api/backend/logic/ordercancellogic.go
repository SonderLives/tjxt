package tradelogic

import (
	"context"
	"errors"

	"tjxt/service/trade/api/backend/model"
	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/auth"
	"tjxt/common/errorx"

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
		return nil, err
	}

	if err := l.svcCtx.OrderModel.UpdateStatus(l.ctx, in.Id, OrderStatusClosed, "用户取消"); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "取消订单失败")
	}

	details, err := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单明细失败")
	}
	for _, d := range details {
		if e := l.svcCtx.OrderDetailModel.UpdateStatus(l.ctx, d.Id, DetailStatusClosed); e != nil {
			return nil, errorx.Wrap(e, errorx.CodeInternal, "关闭订单明细失败")
		}
	}

	return &pb.Empty{}, nil
}
