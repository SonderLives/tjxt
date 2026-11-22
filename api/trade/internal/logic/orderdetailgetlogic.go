// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailGetLogic {
	return &OrderDetailGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailGetLogic) OrderDetailGet(req *types.OrderIdReq) (resp *types.OrderDetailAdminVO, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderDetailGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	resp = &types.OrderDetailAdminVO{
		Id:                 reply.Id,
		OrderId:            reply.OrderId,
		Mobile:             reply.Mobile,
		Name:               reply.Name,
		Price:              reply.Price,
		RealPayAmount:      reply.RealPayAmount,
		DiscountAmount:     reply.DiscountAmount,
		CouponAmount:       reply.CouponAmount,
		CouponDesc:         reply.CouponDesc,
		Status:             int64(reply.Status),
		Message:            reply.Message,
		PayChannel:         reply.PayChannel,
		PayOrderNo:         reply.PayOrderNo,
		StudyValidTime:     reply.StudyValidTime,
		RefundApplyId:      reply.RefundApplyId,
		RefundOrderNo:      reply.RefundOrderNo,
		RefundStatus:       int64(reply.RefundStatus),
		RefundReason:       reply.RefundReason,
		RefundMessage:      reply.RefundMessage,
		RefundChannel:      reply.RefundChannel,
		RefundFailedReason: reply.RefundFailedReason,
		RefundProposerName: reply.RefundProposerName,
		Remark:             reply.Remark,
		CanRefund:          reply.CanRefund,
		FailedReason:       reply.FailedReason,
		Nodes:              make([]types.OrderProgressNodeVO, 0, len(reply.Nodes)),
	}
	for _, n := range reply.Nodes {
		resp.Nodes = append(resp.Nodes, types.OrderProgressNodeVO{
			Name:   n.Name,
			Desc:   n.Desc,
			Status: int64(n.Status),
			Time:   n.Time,
		})
	}
	return resp, nil
}
