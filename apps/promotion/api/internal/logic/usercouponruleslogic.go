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

type UserCouponRulesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponRulesLogic {
	return &UserCouponRulesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponRulesLogic) UserCouponRules(req *types.CouponIdsReq) (resp *types.RulesResp, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.PromotionRpc.UserCouponRules(l.ctx, &pb.IdsRequest{
		Ids:     req.CouponIds,
		UserId:  userId,
		OrderId: req.OrderId,
	})
	if err != nil {
		return nil, err
	}
	return &types.RulesResp{Rules: out.Rules}, nil
}
