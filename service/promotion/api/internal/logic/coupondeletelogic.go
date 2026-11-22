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

type CouponDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponDeleteLogic {
	return &CouponDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponDeleteLogic) CouponDelete(req *types.CouponIdReq) error {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}
	_, err = l.svcCtx.PromotionRpc.CouponDelete(l.ctx, &pb.IdRequest{Id: req.Id, UserId: userId})
	return err
}
