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

type RefundApplyGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyGetLogic {
	return &RefundApplyGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyGetLogic) RefundApplyGet(req *types.RefundIdReq) (resp *types.RefundApplyVO, err error) {
	reply, err := l.svcCtx.TradeRpc.RefundApplyGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.RefundApplyVO{
		Id:             reply.Id,
		OrderId:        reply.OrderId,
		OrderDetailId:  reply.OrderDetailId,
		Price:          reply.Price,
		RefundAmount:   reply.RefundAmount,
		RefundStatus:   int64(reply.RefundStatus),
		RefundOrderNo:  reply.RefundOrderNo,
		PayOrderNo:     reply.PayOrderNo,
		PayChannel:     reply.PayChannel,
		RefundChannel:  reply.RefundChannel,
		RefundReason:   reply.RefundReason,
		RefundMessage:  reply.RefundMessage,
		FailedReason:   reply.FailedReason,
		ApproveOpinion: reply.ApproveOpinion,
		ApproveType:    int64(reply.ApproveType),
		Remark:         reply.Remark,
		CreateTime:     reply.CreateTime,
		OrderTime:      reply.OrderTime,
		PaySuccessTime: reply.PaySuccessTime,
		StudentName:    reply.StudentName,
		Mobile:         reply.Mobile,
	}, nil
}
