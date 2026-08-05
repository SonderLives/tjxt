package logic

import (
	"context"
	"time"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponDiscountLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponDiscountLogic {
	return &UserCouponDiscountLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponDiscount 根据用户选定的券方案，计算订单优惠明细。
// 仅接受属于当前用户且未使用的券，防止越权占用他人优惠券。
func (l *UserCouponDiscountLogic) UserCouponDiscount(in *pb.OrderCouponDTO) (*pb.CouponDiscountDTO, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if len(in.CourseList) == 0 {
		return nil, xerr.BadRequestf("订单课程不能为空")
	}
	if len(in.UserCouponIds) == 0 {
		return &pb.CouponDiscountDTO{}, nil
	}

	userCoupons, err := l.svcCtx.UserCouponModel.FindByIdsAndUser(l.ctx, in.UserCouponIds, in.UserId)
	if err != nil {
		l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "计算优惠明细失败")
	}
	if len(userCoupons) != len(in.UserCouponIds) {
		return nil, xerr.BadRequestf("优惠券不存在或不属于当前用户")
	}

	ids := make([]int64, 0, len(userCoupons))
	for _, uc := range userCoupons {
		ids = append(ids, uc.CouponId)
	}
	list, err := l.svcCtx.CouponModel.FindByIds(l.ctx, ids)
	if err != nil {
		l.Errorf("find coupons by ids failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "计算优惠明细失败")
	}
	coupons := make(map[int64]*model.Coupon, len(list))
	for _, c := range list {
		coupons[c.Id] = c
	}

	available := usableUserCoupons(userCoupons, coupons, time.Now())
	if len(available) != len(in.UserCouponIds) {
		return nil, xerr.Conflict("存在已使用或已过期的优惠券，请重新选择")
	}

	solution, ok := buildSolution(available, in.CourseList)
	if !ok {
		return nil, xerr.Conflict("所选优惠券不满足使用条件")
	}
	return solution, nil
}
