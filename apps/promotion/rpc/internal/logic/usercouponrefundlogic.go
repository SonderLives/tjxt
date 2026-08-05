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

type UserCouponRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponRefundLogic {
	return &UserCouponRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponRefund 退还优惠券，由 trade 服务在订单取消/退款后调用。
// 状态流转 used -> unused；已过期的券即使退还也无法再用，此处仅恢复状态。
func (l *UserCouponRefundLogic) UserCouponRefund(in *pb.IdsRequest) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if len(in.Ids) == 0 {
		return &pb.Empty{}, nil
	}

	userCoupons, err := l.svcCtx.UserCouponModel.FindByIdsAndUser(l.ctx, in.Ids, in.UserId)
	if err != nil {
		l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "退还优惠券失败")
	}
	if len(userCoupons) == 0 {
		return nil, xerr.BadRequestf("优惠券不存在或不属于当前用户")
	}

	rows, err := l.svcCtx.UserCouponModel.UpdateStatusByIds(l.ctx, in.Ids, in.UserId,
		model.UserCouponStatusUsed, model.UserCouponStatusUnused, 0)
	if err != nil {
		l.Errorf("refund user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "退还优惠券失败")
	}
	if rows == 0 {
		return nil, xerr.Conflict("优惠券未被使用，无需退还")
	}

	now := time.Now()
	for _, uc := range userCoupons {
		if uc.Status != model.UserCouponStatusUsed {
			continue
		}
		// 回滚券模板的已使用数量
		if err := l.svcCtx.CouponModel.AddUsedNum(l.ctx, uc.CouponId, -1); err != nil {
			l.Errorf("decr used num failed, couponId=%d, err=%v", uc.CouponId, err)
		}
		if uc.ExpireTime.Valid && now.After(uc.ExpireTime.Time) {
			l.Infof("refunded user coupon already expired, id=%d", uc.Id)
		}
	}
	return &pb.Empty{}, nil
}
