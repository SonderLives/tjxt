package paylogic

import (
	"tjxt/common/timex"
	"tjxt/service/pay/rpc/internal/model"
	"tjxt/service/pay/rpc/pb"
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
		PaySuccessTime: timex.FormatNullTime(m.PaySuccessTime),
		PayOverTime:    timex.Format(m.PayOverTime),
		QrCodeUrl:      formatNullString(m.QrCodeUrl),
		CreateTime:     timex.Format(m.CreateTime),
		UpdateTime:     timex.Format(m.UpdateTime),
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
		CreateTime:        timex.Format(m.CreateTime),
		UpdateTime:        timex.Format(m.UpdateTime),
	}
}
