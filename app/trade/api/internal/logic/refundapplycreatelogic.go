// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/app/trade/api/backend/logic"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/app/trade/api/internal/svc"
	"tjxt/app/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCreateLogic {
	return &RefundApplyCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyCreateLogic) RefundApplyCreate(req *types.RefundFormReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewRefundApplyCreateLogic(l.ctx, l.svcCtx.Backend).RefundApplyCreate(&pb.RefundApplyFormRequest{
		OrderDetailId: req.OrderDetailId,
		RefundReason:  req.RefundReason,
		QuestionDesc:  req.QuestionDesc,
	}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
