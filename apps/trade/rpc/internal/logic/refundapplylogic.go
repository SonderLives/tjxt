package logic

import (
	"context"
	"encoding/json"
	"errors"

	promotionclient "tjxt/apps/promotion/rpc/promotion"
	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/mq"
	"tjxt/pkg/mq/event"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyLogic {
	return &RefundApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 退款 =====
func (l *RefundApplyLogic) RefundApply(in *pb.RefundApplyRequest) (*pb.RefundResultDTO, error) {
	if in.BizOrderNo <= 0 {
		return nil, xerr.BadRequestf("业务订单号不能为空")
	}

	resp, err := l.svcCtx.PayRpc.ApplyRefund(l.ctx, &payclient.ApplyRefundRequest{
		BizOrderNo:       in.BizOrderNo,
		BizRefundOrderNo: in.BizRefundOrderNo,
		RefundAmount:     in.RefundAmount,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "申请退款失败")
	}

	// 退款成功（demo 渠道同步返回成功）→ 通知 learning 撤销课程
	// （best-effort：发布失败仅告警，可由对账/人工兜底）
	if resp.Status == 3 {
		l.publishOrderRefund(in.BizOrderNo)
	}

	return &pb.RefundResultDTO{
		BizPayOrderId:    in.BizOrderNo,
		BizRefundOrderId: resp.BizRefundOrderNo,
		PayOrderNo:       0,
		RefundOrderNo:    resp.RefundOrderNo,
		Status:           resp.Status,
		PayChannel:       "",
		RefundChannel:    "",
	}, nil
}

// publishOrderRefund 向 learning 发布退款成功事件（order.exchange / order.refund），
// 并退还订单使用的优惠券（UserCouponRefund，失败仅告警）。
func (l *RefundApplyLogic) publishOrderRefund(orderId int64) {
	if l.svcCtx.MQProducer == nil {
		logx.Errorf("mq producer unavailable, skip order.refund event, orderId=%d", orderId)
		return
	}
	order, err := l.svcCtx.OrderModel.FindOne(l.ctx, orderId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			logx.Errorf("order not found for order.refund event, orderId=%d", orderId)
			return
		}
		logx.Errorf("query order for order.refund event failed, orderId=%d: %v", orderId, err)
		return
	}

	// 退还订单所用优惠券（幂等：promotion 侧对已退券直接忽略）
	if order.CouponIds.Valid && order.CouponIds.String != "" {
		var couponIds []int64
		if jerr := json.Unmarshal([]byte(order.CouponIds.String), &couponIds); jerr == nil && len(couponIds) > 0 {
			if _, uerr := l.svcCtx.PromotionRpc.UserCouponRefund(l.ctx, &promotionclient.IdsRequest{
				Ids:     couponIds,
				UserId:  order.UserId,
				OrderId: order.Id,
			}); uerr != nil {
				logx.Errorf("refund coupons failed, orderId=%d couponIds=%v: %v", order.Id, couponIds, uerr)
			}
		}
	}

	details, err := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, orderId)
	if err != nil {
		logx.Errorf("list order details for order.refund event failed, orderId=%d: %v", orderId, err)
		return
	}
	courseIds := make([]int64, 0, len(details))
	for _, d := range details {
		courseIds = append(courseIds, d.CourseId)
	}
	if err := l.svcCtx.MQProducer.Publish(l.ctx, mq.ExchangeOrder, mq.RoutingKeyOrderRefund, event.OrderRefundEvent{
		OrderBasic: event.OrderBasic{
			OrderID:    order.Id,
			UserID:     order.UserId,
			CourseIDs:  courseIds,
			FinishTime: now(),
		},
	}); err != nil {
		logx.Errorf("publish order.refund event failed, orderId=%d: %v", orderId, err)
	}
}
