// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyDeleteLogic {
	return &RefundApplyDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyDeleteLogic) RefundApplyDelete(req *types.RefundIdReq) (resp *types.NamePlaceVO, err error) {
	// todo: add your logic here and delete this line

	return
}
