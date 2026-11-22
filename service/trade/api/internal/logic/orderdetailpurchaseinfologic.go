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

type OrderDetailPurchaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailPurchaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPurchaseInfoLogic {
	return &OrderDetailPurchaseInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailPurchaseInfoLogic) OrderDetailPurchaseInfo(req *types.PurchaseInfoReq) (resp *types.PurchaseInfoVO, err error) {
	reply, err := backendlogic.NewOrderDetailPurchaseInfoLogic(l.ctx, l.svcCtx.Backend).OrderDetailPurchaseInfo(&pb.PurchaseInfoRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}

	return &types.PurchaseInfoVO{
		EnrollNum:     reply.EnrollNum,
		RealPayAmount: reply.RealPayAmount,
		RefundNum:     reply.RefundNum,
	}, nil
}
