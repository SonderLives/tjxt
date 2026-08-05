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

type UserCouponExchangeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponExchangeLogic {
	return &UserCouponExchangeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponExchangeLogic) UserCouponExchange(req *types.ExchangeByCodeReq) error {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.PromotionRpc.UserCouponExchange(l.ctx, &pb.ExchangeRequest{
		Code:   req.Code,
		UserId: userId,
	})
	return err
}
