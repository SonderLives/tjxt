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

type UserCouponReceiveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponReceiveLogic {
	return &UserCouponReceiveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponReceiveLogic) UserCouponReceive(req *types.CouponReceiveReq) error {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.PromotionRpc.UserCouponReceive(l.ctx, &pb.IdRequest{
		Id:     req.CouponId,
		UserId: userId,
	})
	return err
}
