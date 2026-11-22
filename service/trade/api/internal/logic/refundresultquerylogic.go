// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundResultQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundResultQueryLogic {
	return &RefundResultQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundResultQueryLogic) RefundResultQuery(req *types.BizRefundOrderIdPathReq) (resp *types.RefundResultDTO, err error) {
	reply, err := backendlogic.NewRefundResultQueryLogic(l.ctx, l.svcCtx.Backend).RefundResultQuery(&pb.RefundResultQueryRequest{
		BizRefundOrderId: req.BizRefundOrderId,
	})
	if err != nil {
		return nil, err
	}
	return &types.RefundResultDTO{
		BizPayOrderId:    reply.BizPayOrderId,
		BizRefundOrderId: reply.BizRefundOrderId,
		PayOrderNo:       reply.PayOrderNo,
		RefundOrderNo:    reply.RefundOrderNo,
		Status:           int64(reply.Status),
		PayChannel:       reply.PayChannel,
		RefundChannel:    reply.RefundChannel,
	}, nil
}
