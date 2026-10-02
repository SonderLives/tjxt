package logic

import (
	"context"
	"time"

	"tjxt/apps/pay/rpc/internal/svc"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/mq"
	"tjxt/pkg/mq/event"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// NotifyPaySuccessLogic 第三方支付渠道/前端模拟回调：标记支付成功。
//
// 真实生产流程：第三方网关 → 我们暴露出去的 notify 接口（带签名校验）→
// 该 RPC 修改本地状态 → 发布 pay.success 事件，trade 消费后回写订单状态。
type NotifyPaySuccessLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewNotifyPaySuccessLogic(ctx context.Context, svcCtx *svc.ServiceContext) *NotifyPaySuccessLogic {
	return &NotifyPaySuccessLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *NotifyPaySuccessLogic) NotifyPaySuccess(in *pb.NotifyPaySuccessRequest) (*pb.EmptyResponse, error) {
	if in.PayOrderNo <= 0 {
		return nil, xerr.BadRequestf("pay_order_no 非法")
	}
	m, err := l.svcCtx.PayOrderModel.FindOneByPayOrderNo(l.ctx, in.PayOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("支付单不存在")
		}
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "查询支付单失败")
	}
	switch m.Status {
	case PayOrderStatusSuccess:
		return &pb.EmptyResponse{}, nil // 幂等
	case PayOrderStatusClosed:
		return nil, xerr.Conflict("订单已关闭，无法标记支付成功")
	case PayOrderStatusPending, PayOrderStatusPaying:
	}

	transit, err := l.svcCtx.PayOrderModel.MarkToSuccess(l.ctx, m.Id, in.ResultCode, in.ResultMsg)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "更新支付单状态失败")
	}
	if !transit {
		// 并发回调下状态已被其他请求流转，按幂等成功处理，不重复发事件
		return &pb.EmptyResponse{}, nil
	}

	// 通知 trade 回写订单状态（best-effort：发布失败仅告警，
	// 一致性由 trade 侧幂等消费 + 对账兜底）
	if l.svcCtx.MQProducer != nil {
		if err := l.svcCtx.MQProducer.Publish(l.ctx, mq.ExchangePay, mq.RoutingKeyPaySuccess, event.PaySuccessEvent{
			BizOrderNo: m.BizOrderNo,
			PayOrderNo: m.PayOrderNo,
			Amount:     m.Amount,
			PayChannel: m.PayChannelCode,
			PayTime:    time.Now(),
		}); err != nil {
			logx.WithContext(l.ctx).Errorf("publish pay.success event failed, bizOrderNo=%d payOrderNo=%d: %v",
				m.BizOrderNo, m.PayOrderNo, err)
		}
	}
	return &pb.EmptyResponse{}, nil
}
