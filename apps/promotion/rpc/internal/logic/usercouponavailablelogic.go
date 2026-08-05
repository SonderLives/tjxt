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

type UserCouponAvailableLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponAvailableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponAvailableLogic {
	return &UserCouponAvailableLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponAvailable 根据订单中的课程，计算当前用户所有可用的优惠券方案，
// 按优惠金额从高到低返回，供下单页选择。
func (l *UserCouponAvailableLogic) UserCouponAvailable(in *pb.OrderCourseListRequest) (*pb.CouponDiscountListReply, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if len(in.CourseList) == 0 {
		return &pb.CouponDiscountListReply{}, nil
	}

	// 1. 取出用户未使用的券
	userCoupons, err := l.svcCtx.UserCouponModel.FindByUserAndStatus(l.ctx, in.UserId, model.UserCouponStatusUnused)
	if err != nil {
		l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询可用优惠券失败")
	}
	if len(userCoupons) == 0 {
		return &pb.CouponDiscountListReply{}, nil
	}

	// 2. 批量加载券规则
	coupons, err := l.loadCoupons(userCoupons)
	if err != nil {
		return nil, err
	}

	// 3. 过滤 + 组合计算
	available := usableUserCoupons(userCoupons, coupons, time.Now())
	solutions := calcSolutions(available, in.CourseList)

	return &pb.CouponDiscountListReply{List: solutions}, nil
}

// loadCoupons 批量加载用户券对应的券规则，避免 N+1 查询。
func (l *UserCouponAvailableLogic) loadCoupons(ucs []*model.UserCoupon) (map[int64]*model.Coupon, error) {
	idSet := make(map[int64]struct{}, len(ucs))
	ids := make([]int64, 0, len(ucs))
	for _, uc := range ucs {
		if _, ok := idSet[uc.CouponId]; ok {
			continue
		}
		idSet[uc.CouponId] = struct{}{}
		ids = append(ids, uc.CouponId)
	}

	list, err := l.svcCtx.CouponModel.FindByIds(l.ctx, ids)
	if err != nil {
		l.Errorf("find coupons by ids failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券规则失败")
	}

	result := make(map[int64]*model.Coupon, len(list))
	for _, c := range list {
		result[c.Id] = c
	}
	return result, nil
}
