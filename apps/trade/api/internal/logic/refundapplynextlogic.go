// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyNextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyNextLogic {
	return &RefundApplyNextLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyNextLogic) RefundApplyNext() (resp *types.RefundApplyVO, err error) {
	// todo: add your logic here and delete this line

	return
}
