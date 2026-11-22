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

type PayOrderApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayOrderApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderApplyLogic {
	return &PayOrderApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayOrderApplyLogic) PayOrderApply(req *types.PayApplyFormDTO) (resp *types.NamePlaceVO, err error) {
	reply, err := backendlogic.NewPayApplyLogic(l.ctx, l.svcCtx.Backend).PayApply(&pb.PayApplyRequest{
		OrderId:        req.OrderId,
		PayChannelCode: req.PayChannelCode,
	})
	if err != nil {
		return nil, err
	}

	return &types.NamePlaceVO{Existed: true, Url: reply.QrUrl, Message: "ok"}, nil
}
