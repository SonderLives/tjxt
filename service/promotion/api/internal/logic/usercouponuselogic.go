// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/api/promotion/internal/svc"
	"tjxt/api/promotion/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/rpc/promotion/client/promotion"
)

type UserCouponUseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponUseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponUseLogic {
	return &UserCouponUseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponUseLogic) UserCouponUse(req *types.CouponIdsReq) error {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.PromotionRpc.UserCouponUse(l.ctx, &pb.IdsRequest{
		Ids:     req.CouponIds,
		UserId:  userId,
		OrderId: req.OrderId,
	})
	return err
}
