package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponUseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponUseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponUseLogic {
	return &UserCouponUseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponUse 核销优惠券，由 trade 服务在下单成功后调用。
// 状态流转 unused -> used，条件更新保证同一张券不会被重复核销。
func (l *UserCouponUseLogic) UserCouponUse(in *pb.IdsRequest) (*pb.Empty, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}
	if len(in.Ids) == 0 {
		return &pb.Empty{}, nil
	}

	userCoupons, err := l.svcCtx.UserCouponModel.FindByIdsAndUser(l.ctx, in.Ids, in.UserId)
	if err != nil {
		l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "核销优惠券失败")
	}
	if len(userCoupons) != len(in.Ids) {
		return nil, xerr.BadRequestf("优惠券不存在或不属于当前用户")
	}

	rows, err := l.svcCtx.UserCouponModel.UpdateStatusByIds(l.ctx, in.Ids, in.UserId,
		model.UserCouponStatusUnused, model.UserCouponStatusUsed, in.OrderId)
	if err != nil {
		l.Errorf("use user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "核销优惠券失败")
	}
	if rows != int64(len(in.Ids)) {
		return nil, xerr.Conflict("部分优惠券已被使用或已过期")
	}

	// 同步券模板的已使用数量，失败不阻断主流程（统计类字段，可离线校准）
	for _, uc := range userCoupons {
		if err := l.svcCtx.CouponModel.AddUsedNum(l.ctx, uc.CouponId, 1); err != nil {
			l.Errorf("incr used num failed, couponId=%d, err=%v", uc.CouponId, err)
		}
	}
	return &pb.Empty{}, nil
}
