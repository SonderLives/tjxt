package logic

import (
	"context"
	"time"

	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponListLogic {
	return &CouponListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponList 查询正在发放中的优惠券列表（C 端），
// 并标记每张券当前用户是否可领（available）与是否已领（received）。
func (l *CouponListLogic) CouponList(in *pb.CouponListRequest) (*pb.CouponListReply, error) {
	list, err := l.svcCtx.CouponModel.FindList(l.ctx)
	if err != nil {
		l.Errorf("list coupon failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券失败")
	}

	// 一次性拉取当前用户已领券，避免逐张券查库
	received := make(map[int64]int64)
	if in.UserId > 0 {
		mine, err := l.svcCtx.UserCouponModel.FindByUserAndStatus(l.ctx, in.UserId, "")
		if err != nil {
			l.Errorf("find user coupon failed, userId=%d, err=%v", in.UserId, err)
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券失败")
		}
		for _, uc := range mine {
			received[uc.CouponId]++
		}
	}

	now := time.Now()
	items := make([]*pb.CouponVO, 0, len(list))
	for _, c := range list {
		if !couponReceivable(c, now) {
			// 仅展示发放中的券
			continue
		}
		got := received[c.Id]
		// 达到限领数量后不可再领
		available := c.UserLimit == 0 || got < c.UserLimit
		items = append(items, toCouponVO(c, available, got > 0))
	}
	return &pb.CouponListReply{List: items}, nil
}
