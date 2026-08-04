package logic

import (
	"context"
	"database/sql"

	"tjxt/apps/pay/rpc/internal/model"
	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// ApplyRefundLogic 申请退款（创建 or 幂等返回已有退款单）
type ApplyRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewApplyRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyRefundLogic {
	return &ApplyRefundLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *ApplyRefundLogic) ApplyRefund(in *pb.ApplyRefundRequest) (*pb.RefundResultResponse, error) {
	if in.BizOrderNo <= 0 || in.BizRefundOrderNo <= 0 || in.RefundAmount <= 0 {
		return nil, xerr.BadRequestf("biz_order_no/biz_refund_order_no/refund_amount 非法")
	}

	// 幂等：biz_refund_order_no 已有退款则直接返回
	if existing, err := l.svcCtx.RefundOrderModel.FindOneByBizRefundOrderNo(l.ctx, in.BizRefundOrderNo); err == nil {
		return &pb.RefundResultResponse{
			RefundOrderNo:    existing.RefundOrderNo,
			BizRefundOrderNo: existing.BizRefundOrderNo,
			RefundAmount:     existing.RefundAmount,
			Status:           int32(existing.Status),
			ResultMsg:        existing.ResultMsg,
		}, nil
	} else if !isNotFound(err) {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询已有退款单失败")
	}

	// 校验原支付单已成功
	payOrder, err := l.svcCtx.PayOrderModel.FindOneByBizOrderNo(l.ctx, in.BizOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("原支付单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询原支付单失败")
	}
	if payOrder.Status != PayOrderStatusSuccess {
		return nil, xerr.Conflict("原支付单未支付成功，无法退款")
	}
	if in.RefundAmount > payOrder.Amount {
		return nil, xerr.BadRequestf("退款金额超过原支付金额")
	}
	// 累计已退款金额校验：防止多次拆单退款超额
	refunds, err := l.svcCtx.RefundOrderModel.FindListByBizOrderNo(l.ctx, in.BizOrderNo)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询已有退款记录失败")
	}
	var refunded int64
	for _, r := range refunds {
		if r.Status == RefundStatusSuccess || r.Status == RefundStatusProcessing {
			refunded += r.RefundAmount
		}
	}
	if refunded+in.RefundAmount > payOrder.Amount {
		return nil, xerr.BadRequestf("累计退款金额超过原支付金额")
	}

	roNo := nextID()
	ro := &model.RefundOrder{
		BizOrderNo:       in.BizOrderNo,
		BizRefundOrderNo: in.BizRefundOrderNo,
		PayOrderNo:       payOrder.PayOrderNo,
		RefundOrderNo:    roNo,
		RefundAmount:     in.RefundAmount,
		TotalAmount:      payOrder.Amount,
		IsSplit:          0,
		PayChannelCode:   payOrder.PayChannelCode,
		Status:           RefundStatusProcessing,
		NotifyStatus:     RefundNotifyStatusPending,
	}
	if _, err := l.svcCtx.RefundOrderModel.Insert(l.ctx, ro); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "创建退款单失败")
	}

	// 真实生产：调用第三方退款 API，再根据结果 async 通过 NotifyRefundSuccess/Failed 更新
	// demo 中暂时直接 mock 成功
	_ = sql.ErrNoRows
	if err := l.svcCtx.RefundOrderModel.MarkToSuccess(l.ctx, ro.Id, "MOCK_OK", "mock 退款成功", "mock"); err != nil {
		l.Errorf("mock 退款成功更新失败: %v", err)
	}

	return &pb.RefundResultResponse{
		RefundOrderNo:    roNo,
		BizRefundOrderNo: in.BizRefundOrderNo,
		RefundAmount:     in.RefundAmount,
		Status:           RefundStatusSuccess,
		ResultMsg:        "mock 退款成功",
	}, nil
}