package logic

import (
	"database/sql"
	"time"

	"tjxt/apps/pay/rpc/internal/model"
	"tjxt/apps/pay/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/utils/page"
)

// 支付渠道状态
const (
	PayChannelStatusEnabled  = 1 // 使用中
	PayChannelStatusDisabled = 2 // 停用
)

// 支付单状态（与表注释一致）
const (
	PayOrderStatusPending = 0 // 待提交
	PayOrderStatusPaying  = 1 // 待支付
	PayOrderStatusClosed  = 2 // 支付超时或取消
	PayOrderStatusSuccess = 3 // 支付成功
)

// 退款单状态
const (
	RefundStatusInit       = 0 // 未提交
	RefundStatusProcessing = 1 // 退款中
	RefundStatusFailed     = 2 // 退款失败
	RefundStatusSuccess    = 3 // 退款成功
)

// 回调状态（PayOrder.NotifyStatus）
const (
	NotifyStatusPending = 0 // 待回调
	NotifyStatusOK      = 1 // 回调成功
	NotifyStatusFail    = 2 // 回调失败
)

// 退款通知状态（RefundOrder.NotifyStatus）
const (
	RefundNotifyStatusPending    = 0 // 待通知
	RefundNotifyStatusSuccess    = 1 // 通知成功
	RefundNotifyStatusProcessing = 2 // 通知中
	RefundNotifyStatusFailed     = 3 // 通知失败
)

// 渠道类型（与表注释一致，1=h5 2=小程序 3=公众号 4=扫码）
const (
	PayTypeH5     = 1
	PayTypeMini   = 2
	PayTypeMp     = 3
	PayTypeNative = 4
)

func toPayChannelResp(m *model.PayChannel) *pb.PayChannelResponse {
	return &pb.PayChannelResponse{
		Id:              m.Id,
		Name:            m.Name,
		ChannelCode:     m.ChannelCode,
		ChannelPriority: int32(m.ChannelPriority),
		ChannelIcon:     m.ChannelIcon,
		Status:          int32(m.Status),
	}
}

func toPayOrderResp(m *model.PayOrder) *pb.PayOrderResponse {
	return &pb.PayOrderResponse{
		Id:             m.Id,
		BizOrderNo:     m.BizOrderNo,
		PayOrderNo:     m.PayOrderNo,
		BizUserId:      m.BizUserId,
		PayChannelCode: m.PayChannelCode,
		Amount:         m.Amount,
		PayType:        int32(m.PayType),
		Status:         int32(m.Status),
		ExpandJson:     m.ExpandJson,
		NotifyUrl:      m.NotifyUrl,
		NotifyTimes:    int32(m.NotifyTimes),
		NotifyStatus:   int32(m.NotifyStatus),
		ResultCode:     m.ResultCode,
		ResultMsg:      m.ResultMsg,
		PaySuccessTime: formatNullTime(m.PaySuccessTime),
		PayOverTime:    formatTime(m.PayOverTime),
		QrCodeUrl:      formatNullString(m.QrCodeUrl),
		CreateTime:     formatTime(m.CreateTime),
		UpdateTime:     formatTime(m.UpdateTime),
	}
}

func toRefundResp(m *model.RefundOrder) *pb.RefundOrderResponse {
	return &pb.RefundOrderResponse{
		Id:                m.Id,
		BizOrderNo:        m.BizOrderNo,
		BizRefundOrderNo:  m.BizRefundOrderNo,
		PayOrderNo:        m.PayOrderNo,
		RefundOrderNo:     m.RefundOrderNo,
		RefundAmount:      m.RefundAmount,
		TotalAmount:       m.TotalAmount,
		IsSplit:           m.IsSplit == 1,
		PayChannelCode:    m.PayChannelCode,
		ResultCode:        m.ResultCode,
		ResultMsg:         m.ResultMsg,
		Status:            int32(m.Status),
		RefundChannel:     formatNullString(m.RefundChannel),
		NotifyFailedTimes: int32(m.NotifyFailedTimes),
		NotifyStatus:      int32(m.NotifyStatus),
		CreateTime:        formatTime(m.CreateTime),
		UpdateTime:        formatTime(m.UpdateTime),
	}
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format("2006-01-02 15:04:05")
}

func formatNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// nextID 雪花算法生成分布式 ID
func nextID() int64 { return idgen.NextID() }

// isNotFound 判断错误是否是数据不存在
func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// normalizePage 把 proto 的分页参数归一化成 offset/limit
func normalizePage(p *pb.PageQueryPayChannelsRequest) (offset, limit int64) {
	return page.Normalize(p.PageNo, p.PageSize)
}

// buildPayChannel 把 proto 请求转换为 model 数据结构（供新增/更新复核用）
func buildPayChannel(in *pb.PayChannelRequest) *model.PayChannel {
	return &model.PayChannel{
		Id:              in.Id,
		Name:            in.Name,
		ChannelCode:     in.ChannelCode,
		ChannelPriority: int64(in.ChannelPriority),
		ChannelIcon:     in.ChannelIcon,
		Status:          PayChannelStatusEnabled,
	}
}