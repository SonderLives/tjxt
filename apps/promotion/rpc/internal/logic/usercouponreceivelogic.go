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

type UserCouponReceiveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponReceiveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponReceiveLogic {
	return &UserCouponReceiveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponReceive 领取优惠券。
// 库存扣减依赖 SQL 条件更新保证并发安全：先抢库存再写用户券，
// 写用户券失败时回滚已扣减的库存。
func (l *UserCouponReceiveLogic) UserCouponReceive(in *pb.IdRequest) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("优惠券 id 非法")
	}

	coupon, err := l.svcCtx.CouponModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("优惠券不存在")
		}
		l.Errorf("find coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "领取优惠券失败")
	}
	if coupon.Deleted == 1 {
		return nil, xerr.NotFound("优惠券不存在")
	}
	if coupon.ObtainWay != model.ObtainWayReceive {
		return nil, xerr.Conflict("该优惠券不支持手动领取")
	}

	now := time.Now()
	if !couponReceivable(coupon, now) {
		return nil, xerr.Conflict("优惠券不在领取时间内或已被领完")
	}

	// 限领校验（非强一致，最终由库存与业务容忍度兜底）
	if coupon.UserLimit > 0 {
		got, err := l.svcCtx.UserCouponModel.CountByUserAndCoupon(l.ctx, in.UserId, coupon.Id)
		if err != nil {
			l.Errorf("count user coupon failed, userId=%d, couponId=%d, err=%v", in.UserId, coupon.Id, err)
			return nil, xerr.Wrap(err, xerr.CodeInternal, "领取优惠券失败")
		}
		if got >= coupon.UserLimit {
			return nil, xerr.Conflict("已达到该优惠券的领取上限")
		}
	}

	// 原子扣减库存
	rows, err := l.svcCtx.CouponModel.IncrIssueNum(l.ctx, coupon.Id)
	if err != nil {
		l.Errorf("incr issue num failed, couponId=%d, err=%v", coupon.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "领取优惠券失败")
	}
	if rows == 0 {
		return nil, xerr.Conflict("优惠券已被领完")
	}

	uc := &model.UserCoupon{
		UserId:     in.UserId,
		CouponId:   coupon.Id,
		Status:     model.UserCouponStatusUnused,
		ObtainTime: sqlNow(),
		ExpireTime: userCouponExpireTime(coupon, now),
		Creater:    in.UserId,
		Updater:    in.UserId,
	}
	if _, err := l.svcCtx.UserCouponModel.Insert(l.ctx, uc); err != nil {
		// 用户券落库失败，回滚已扣减的库存，避免券被“吞掉”
		if rbErr := l.svcCtx.CouponModel.DecrIssueNum(l.ctx, coupon.Id); rbErr != nil {
			l.Errorf("rollback issue num failed, couponId=%d, err=%v", coupon.Id, rbErr)
		}
		l.Errorf("insert user coupon failed, userId=%d, couponId=%d, err=%v", in.UserId, coupon.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "领取优惠券失败")
	}

	return &pb.Empty{}, nil
}
