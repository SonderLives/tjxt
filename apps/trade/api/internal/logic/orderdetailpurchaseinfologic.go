// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
