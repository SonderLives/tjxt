// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyPageLogic {
	return &RefundApplyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyPageLogic) RefundApplyPage(req *types.RefundApplyPageReq) (resp *types.RefundApplyPageReply, err error) {
	// todo: add your logic here and delete this line

	return
}
