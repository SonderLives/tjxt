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

type OrderDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDeleteLogic {
	return &OrderDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDeleteLogic) OrderDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}

	// 先删明细，再删主单
	details, err := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单明细失败")
	}
	for _, d := range details {
		if e := l.svcCtx.OrderDetailModel.Delete(l.ctx, d.Id); e != nil {
			return nil, errorx.Wrap(e, errorx.CodeInternal, "删除订单明细失败")
		}
	}

	if err = l.svcCtx.OrderModel.Delete(l.ctx, in.Id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "删除订单失败")
	}

	return &pb.Empty{}, nil
}
