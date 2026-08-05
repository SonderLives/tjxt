// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/promotion/api/internal/svc"
	"tjxt/apps/promotion/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/apps/promotion/rpc/promotion"
)

type CouponGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponGetLogic {
	return &CouponGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponGetLogic) CouponGet(req *types.CouponIdReq) (resp *types.CouponDetailVO, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.PromotionRpc.CouponGet(l.ctx, &pb.IdRequest{Id: req.Id, UserId: userId})
	if err != nil {
		return nil, err
	}
	return toCouponDetailVO(out), nil
}
