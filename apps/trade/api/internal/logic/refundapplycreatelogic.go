// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
