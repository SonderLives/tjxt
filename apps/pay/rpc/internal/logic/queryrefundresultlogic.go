package logic

import (
	"context"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryRefundResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryRefundResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryRefundResultLogic {
	return &QueryRefundResultLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QueryRefundResultLogic) QueryRefundResult(in *pb.QueryRefundResultRequest) (*pb.RefundResultResponse, error) {
	if in.BizRefundOrderNo <= 0 {
		return nil, xerr.BadRequestf("biz_refund_order_no 非法")
	}
	m, err := l.svcCtx.RefundOrderModel.FindOneByBizRefundOrderNo(l.ctx, in.BizRefundOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("退款单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询退款结果失败")
	}
	return &pb.RefundResultResponse{
		RefundOrderNo:    m.RefundOrderNo,
		BizRefundOrderNo: m.BizRefundOrderNo,
		RefundAmount:     m.RefundAmount,
		Status:           int32(m.Status),
		ResultMsg:        m.ResultMsg,
	}, nil
}