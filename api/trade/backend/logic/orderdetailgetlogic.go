package tradelogic

import (
	"context"
	"errors"

	"tjxt/api/trade/backend/model"
	"tjxt/api/trade/backend/svc"
	"tjxt/api/trade/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailGetLogic {
	return &OrderDetailGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 订单明细 =====
func (l *OrderDetailGetLogic) OrderDetailGet(in *pb.IdRequest) (*pb.OrderDetailAdminVO, error) {
	d, err := l.svcCtx.OrderDetailModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单明细不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单明细失败")
	}

	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, d.OrderId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单失败")
	}

	// 退款申请可能不存在，忽略未找到
	ra, err := l.svcCtx.RefundApplyModel.FindByOrderDetailId(l.ctx, d.Id)
	if err != nil {
		if !errors.Is(err, model.ErrNotFound) {
			return nil, errorx.Wrap(err, errorx.CodeInternal, "查询退款申请失败")
		}
		ra = nil
	}

	return toOrderDetailAdminVO(d, order, ra), nil
}
