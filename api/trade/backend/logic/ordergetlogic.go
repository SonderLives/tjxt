package tradelogic

import (
	"context"
	"errors"

	"tjxt/api/trade/backend/model"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

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
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单失败")
	}

	details, e := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if e != nil && !errors.Is(e, model.ErrNotFound) {
		return nil, errorx.Wrap(e, errorx.CodeInternal, "查询订单明细失败")
	}

	return toOrderVO(order, details), nil
}
