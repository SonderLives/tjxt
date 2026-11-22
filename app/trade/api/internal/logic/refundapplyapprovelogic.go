// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyApproveLogic {
	return &RefundApplyApproveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyApproveLogic) RefundApplyApprove(req *types.ApproveFormReq) (resp *types.NamePlaceVO, err error) {
	if _, err = backendlogic.NewRefundApplyApproveLogic(l.ctx, l.svcCtx.Backend).RefundApplyApprove(&pb.ApproveRequest{
		Id:             req.Id,
		ApproveType:    req.ApproveType,
		ApproveOpinion: req.ApproveOpinion,
		Remark:         req.Remark,
	}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
