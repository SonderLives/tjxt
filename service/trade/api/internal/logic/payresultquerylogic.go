// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/api/trade/backend/logic"
	"tjxt/api/trade/backend/pb"
	"tjxt/api/trade/internal/svc"
	"tjxt/api/trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayResultQueryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayResultQueryLogic {
	return &PayResultQueryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayResultQueryLogic) PayResultQuery(req *types.BizOrderIdPathReq) (resp *types.PayResultDTO, err error) {
	reply, err := backendlogic.NewPayResultQueryLogic(l.ctx, l.svcCtx.Backend).PayResultQuery(&pb.PayResultQueryRequest{BizOrderId: req.BizOrderId})
	if err != nil {
		return nil, err
	}

	return &types.PayResultDTO{
		BizOrderId:  reply.BizOrderId,
		Status:      int64(reply.Status),
		PayChannel:  reply.PayChannel,
		PayOrderNo:  reply.PayOrderNo,
		SuccessTime: reply.SuccessTime,
	}, nil
}
