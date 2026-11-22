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
	if _, err = backendlogic.NewRefundApplyDeleteLogic(l.ctx, l.svcCtx.Backend).RefundApplyDelete(&pb.IdRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
