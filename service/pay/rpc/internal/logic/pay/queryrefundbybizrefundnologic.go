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

// QueryRefundByBizRefundNoLogic 按业务退款单号查询退款详情
type QueryRefundByBizRefundNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRefundByBizRefundNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRefundByBizRefundNoLogic {
	return &QueryRefundByBizRefundNoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QueryRefundByBizRefundNoLogic) QueryRefundByBizRefundNo(in *pb.QueryRefundRequest) (*pb.RefundOrderResponse, error) {
	if in.BizRefundOrderNo <= 0 {
		return nil, errorx.BadRequestf("biz_refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByBizRefundOrderNo(l.ctx, in.BizRefundOrderNo)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("退款单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询退款单失败")
	}
	return toRefundResp(m), nil
}
