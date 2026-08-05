package logic

import (
	"context"
	"strings"
	"time"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponExchangeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponExchangeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponExchangeLogic {
	return &UserCouponExchangeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponExchange 使用兑换码兑换优惠券。
// 兑换码核销与库存扣减均通过条件更新保证并发安全，任一环节失败都会补偿前一步。
func (l *UserCouponExchangeLogic) UserCouponExchange(in *pb.ExchangeRequest) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	code := strings.ToUpper(strings.TrimSpace(in.Code))
	if code == "" {
		return nil, xerr.BadRequestf("兑换码不能为空")
	}

	// 1. 校验兑换码
	cc, err := l.svcCtx.CouponCodeModel.FindOneByCode(l.ctx, code)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("兑换码不存在")
		}
		l.Errorf("find coupon code failed, code=%s, err=%v", code, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "兑换失败")
	}
	if cc.Deleted == 1 {
		return nil, xerr.NotFound("兑换码不存在")
	}
	if cc.Status != model.CouponCodeStatusUnused {
		return nil, xerr.Conflict("兑换码已被使用")
	}

	now := time.Now()
	if cc.ExpireTime.Valid && now.After(cc.ExpireTime.Time) {
		return nil, xerr.Conflict("兑换码已过期")
	}

	// 2. 校验优惠券
	coupon, err := l.svcCtx.CouponModel.FindOne(l.ctx, cc.CouponId)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("优惠券不存在")
		}
		l.Errorf("find coupon failed, id=%d, err=%v", cc.CouponId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "兑换失败")
	}
	if coupon.Deleted == 1 || coupon.Status == model.CouponStatusEnded {
		return nil, xerr.Conflict("优惠券活动已结束")
	}
	if coupon.UserLimit > 0 {
		got, err := l.svcCtx.UserCouponModel.CountByUserAndCoupon(l.ctx, in.UserId, coupon.Id)
		if err != nil {
			l.Errorf("count user coupon failed, userId=%d, couponId=%d, err=%v", in.UserId, coupon.Id, err)
			return nil, xerr.Wrap(err, xerr.CodeInternal, "兑换失败")
		}
		if got >= coupon.UserLimit {
			return nil, xerr.Conflict("已达到该优惠券的领取上限")
		}
	}

	// 3. 核销兑换码（条件更新，防并发抢兑）
	rows, err := l.svcCtx.CouponCodeModel.MarkUsed(l.ctx, cc.Id, in.UserId)
	if err != nil {
		l.Errorf("mark coupon code used failed, id=%d, err=%v", cc.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "兑换失败")
	}
	if rows == 0 {
		return nil, xerr.Conflict("兑换码已被使用")
	}

	// 4. 扣减库存并发放用户券
	if _, err := l.svcCtx.CouponModel.IncrIssueNum(l.ctx, coupon.Id); err != nil {
		l.Errorf("incr issue num failed, couponId=%d, err=%v", coupon.Id, err)
	}

	uc := &model.UserCoupon{
		UserId:     in.UserId,
		CouponId:   coupon.Id,
		Status:     model.UserCouponStatusUnused,
		ObtainTime: sqlNow(),
		ExpireTime: userCouponExpireTime(coupon, now),
		Code:       code,
		Creater:    in.UserId,
		Updater:    in.UserId,
	}
	if _, err := l.svcCtx.UserCouponModel.Insert(l.ctx, uc); err != nil {
		l.Errorf("insert user coupon failed, userId=%d, couponId=%d, err=%v", in.UserId, coupon.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "兑换失败")
	}

	return &pb.Empty{}, nil
}
