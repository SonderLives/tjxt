package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponRulesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponRulesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponRulesLogic {
	return &UserCouponRulesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponRules 查询指定用户券的规则文案，用于订单页展示"已优惠"明细。
func (l *UserCouponRulesLogic) UserCouponRules(in *pb.IdsRequest) (*pb.RulesReply, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if len(in.Ids) == 0 {
		return &pb.RulesReply{}, nil
	}

	userCoupons, err := l.svcCtx.UserCouponModel.FindByIdsAndUser(l.ctx, in.Ids, in.UserId)
	if err != nil {
		l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券规则失败")
	}
	if len(userCoupons) == 0 {
		return &pb.RulesReply{}, nil
	}

	couponIds := make([]int64, 0, len(userCoupons))
	for _, uc := range userCoupons {
		couponIds = append(couponIds, uc.CouponId)
	}
	list, err := l.svcCtx.CouponModel.FindByIds(l.ctx, couponIds)
	if err != nil {
		l.Errorf("find coupons by ids failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券规则失败")
	}
	coupons := make(map[int64]*model.Coupon, len(list))
	for _, c := range list {
		coupons[c.Id] = c
	}

	// 保持与入参 ids 一致的顺序返回
	rules := make([]string, 0, len(userCoupons))
	for _, uc := range userCoupons {
		if c, ok := coupons[uc.CouponId]; ok {
			rules = append(rules, couponRule(c))
		}
	}
	return &pb.RulesReply{Rules: rules}, nil
}
