// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
