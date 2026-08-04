package service

import (
	"context"
	"database/sql"
	"strconv"
	"time"

	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"
	"tjxt/apps/trade/api/internal/model"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// PayService 支付业务接口
type PayService interface {
	// ApplyPay 发起支付（模拟支付成功后发布开课事件）。
	ApplyPay(ctx context.Context, userId, orderId int64, channelCode string) (string, error)
	// ListChannels 查询可用支付渠道。
	ListChannels(ctx context.Context) ([]types.PayChannelVO, error)
}

type payService struct {
	orderModel     *model.OrderModel
	detailModel    *model.OrderDetailModel
	channelModel   *model.PayChannelModel
	eventPublisher EventPublisher
}

// NewPayService 创建支付业务服务。
func NewPayService(orderModel *model.OrderModel, detailModel *model.OrderDetailModel, channelModel *model.PayChannelModel, publisher EventPublisher) PayService {
	return &payService{
		orderModel:     orderModel,
		detailModel:    detailModel,
		channelModel:   channelModel,
		eventPublisher: publisher,
	}
}

// ApplyPay 发起支付。
//
// 说明：真实环境中支付流程由 pay 服务承接（申请支付单 -> 扫码 -> 异步回调）。
// 本实现为可运行的闭环模拟：校验订单后直接标记支付成功，并发布订单支付事件，
// 供 learning 等服务通过 MQ 完成课程开通。接入真实支付渠道时替换本方法即可。
func (s *payService) ApplyPay(ctx context.Context, userId, orderId int64, channelCode string) (string, error) {
	if orderId == 0 || channelCode == "" {
		return "", xerr.BadRequestf("订单与支付渠道不能为空")
	}

	channel, err := s.channelModel.FindByCode(ctx, channelCode)
	if err == sql.ErrNoRows {
		return "", xerr.BadRequestf("支付渠道不存在或已停用")
	}
	if err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "查询支付渠道失败")
	}

	order, err := s.orderModel.FindById(ctx, orderId)
	if err == sql.ErrNoRows {
		return "", xerr.NotFound("订单不存在")
	}
	if err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}
	if order.UserId != userId {
		return "", xerr.Forbidden("无权操作该订单")
	}
	if order.Status != model.OrderStatusPendingPay {
		return "", xerr.Conflict("订单状态不允许支付")
	}

	details, err := s.detailModel.ListByOrderId(ctx, order.Id)
	if err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}
	courseIds := make([]int64, 0, len(details))
	now := time.Now()
	for i := range details {
		d := &details[i]
		if d.Status != model.OrderDetailStatusPendingPay {
			continue
		}
		courseIds = append(courseIds, d.CourseId)
		expireTime := expireTimeFrom(d.ValidDuration, now)
		if err := s.detailModel.MarkPaid(ctx, d.Id, model.OrderDetailStatusPaid, channel.ChannelCode, &expireTime); err != nil {
			return "", xerr.Wrap(err, xerr.CodeInternal, "更新订单明细失败")
		}
	}
	if len(courseIds) == 0 {
		return "", xerr.Conflict("订单中没有待支付课程")
	}

	payOrderNo := idgen.NextID()
	finishTime := now.Add(30 * 24 * time.Hour) // 支付后 30 天订单完成
	if err := s.orderModel.MarkPaid(ctx, order.Id, payOrderNo, channel.ChannelCode, finishTime); err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "更新订单失败")
	}

	// 发布支付成功事件，learning 服务据此为用户开通课程
	if err := s.eventPublisher.PublishPay(ctx, order.Id, userId, courseIds, now); err != nil {
		logx.Errorf("publish pay event failed, order=%d err=%v", order.Id, err)
		return "", err
	}

	// 模拟支付返回：真实场景为支付二维码链接 / 唤起参数
	return simulatedPayURL(order.Id, payOrderNo), nil
}

// ListChannels 查询可用支付渠道。
func (s *payService) ListChannels(ctx context.Context) ([]types.PayChannelVO, error) {
	rows, err := s.channelModel.ListEnabled(ctx)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询支付渠道失败")
	}
	list := make([]types.PayChannelVO, 0, len(rows))
	for i := range rows {
		list = append(list, types.PayChannelVO{
			Id:              rows[i].Id,
			Name:            rows[i].Name,
			ChannelCode:     rows[i].ChannelCode,
			ChannelPriority: rows[i].ChannelPriority,
			ChannelIcon:     rows[i].ChannelIcon,
		})
	}
	return list, nil
}

// expireTimeFrom 根据有效期（月）计算课程过期时间。
func expireTimeFrom(validDuration sql.NullInt64, base time.Time) time.Time {
	if validDuration.Valid && validDuration.Int64 > 0 {
		return base.AddDate(0, int(validDuration.Int64), 0)
	}
	// 未配置有效期按 30 天学习期处理，可结合实际需求调整
	return base.Add(30 * 24 * time.Hour)
}

// simulatedPayURL 构造模拟支付返回。
func simulatedPayURL(orderId, payOrderNo int64) string {
	return "simulated://pay?orderId=" + strconv.FormatInt(orderId, 10) + "&payOrderNo=" + strconv.FormatInt(payOrderNo, 10)
}
