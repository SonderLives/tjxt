// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyGetLogic {
	return &RefundApplyGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyGetLogic) RefundApplyGet(req *types.RefundIdReq) (resp *types.RefundApplyVO, err error) {
	// todo: add your logic here and delete this line

	return
}
