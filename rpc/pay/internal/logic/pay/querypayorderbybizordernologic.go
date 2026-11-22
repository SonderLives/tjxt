package paylogic

import (
	"context"

	"tjxt/rpc/pay/internal/svc"
	"tjxt/rpc/pay/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPayOrderByBizOrderNoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryPayOrderByBizOrderNoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPayOrderByBizOrderNoLogic {
	return &QueryPayOrderByBizOrderNoLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QueryPayOrderByBizOrderNoLogic) QueryPayOrderByBizOrderNo(in *pb.QueryPayOrderRequest) (*pb.PayOrderResponse, error) {
	if in.BizOrderNo <= 0 {
		return nil, errorx.BadRequestf("biz_order_no 非法")
	}
	m, err := l.svcCtx.PayOrderModel.FindOneByBizOrderNo(l.ctx, in.BizOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("支付单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付单失败")
	}
	return toPayOrderResp(m), nil
}
