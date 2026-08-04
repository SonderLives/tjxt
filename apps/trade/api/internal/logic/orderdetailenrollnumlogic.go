// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailEnrollNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailEnrollNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailEnrollNumLogic {
	return &OrderDetailEnrollNumLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailEnrollNumLogic) OrderDetailEnrollNum(req *types.EnrollNumReq) (resp *types.NamePlaceVO, err error) {
	// todo: add your logic here and delete this line

	return
}
