package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderGetLogic {
	return &OrderGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderGetLogic) OrderGet(in *pb.IdRequest) (*pb.OrderVO, error) {
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}

	details, e := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if e != nil && !errors.Is(e, model.ErrNotFound) {
		return nil, xerr.Wrap(e, xerr.CodeInternal, "查询订单明细失败")
	}

	return toOrderVO(order, details), nil
}
