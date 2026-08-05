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
	reply, err := l.svcCtx.TradeRpc.OrderDetailPurchaseInfo(l.ctx, &pb.PurchaseInfoRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}

	return &types.PurchaseInfoVO{
		EnrollNum:     reply.EnrollNum,
		RealPayAmount: reply.RealPayAmount,
		RefundNum:     reply.RefundNum,
	}, nil
}
