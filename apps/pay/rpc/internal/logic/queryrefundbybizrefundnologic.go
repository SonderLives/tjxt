package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("biz_refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByBizRefundOrderNo(l.ctx, in.BizRefundOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("退款单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询退款单失败")
	}
	return toRefundResp(m), nil
}