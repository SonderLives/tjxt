package logic

import (
	"context"
	"time"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserCouponPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponPageLogic {
	return &UserCouponPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UserCouponPage 分页查询我的优惠券。
func (l *UserCouponPageLogic) UserCouponPage(in *pb.UserCouponPageRequest) (*pb.UserCouponPageReply, error) {
	if in.UserId <= 0 {
		return nil, xerr.Unauthorized("")
	}

	var pageNo, pageSize int64
	if in.Page != nil {
		pageNo, pageSize = in.Page.PageNo, in.Page.PageSize
	}
	offset, limit := page.Normalize(pageNo, pageSize)

	userCoupons, total, err := l.svcCtx.UserCouponModel.FindPageByUser(l.ctx, in.UserId, in.Status, offset, limit)
	if err != nil {
		l.Errorf("page user coupon failed, userId=%d, err=%v", in.UserId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询我的优惠券失败")
	}
	if len(userCoupons) == 0 {
		return &pb.UserCouponPageReply{Total: total}, nil
	}

	// 批量加载券规则，避免逐条查询
	ids := make([]int64, 0, len(userCoupons))
	for _, uc := range userCoupons {
		ids = append(ids, uc.CouponId)
	}
	list, err := l.svcCtx.CouponModel.FindByIds(l.ctx, ids)
	if err != nil {
		l.Errorf("find coupons by ids failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询我的优惠券失败")
	}
	coupons := make(map[int64]*model.Coupon, len(list))
	for _, c := range list {
		coupons[c.Id] = c
	}

	now := time.Now()
	items := make([]*pb.CouponVO, 0, len(userCoupons))
	for _, uc := range userCoupons {
		c, ok := coupons[uc.CouponId]
		if !ok {
			continue
		}
		vo := toCouponVO(c, uc.Status == model.UserCouponStatusUnused, true)
		// 我的券展示实际到期时间，而非券模板上的有效期
		if uc.ExpireTime.Valid {
			vo.TermEndTime = formatNullTime(uc.ExpireTime)
			if now.After(uc.ExpireTime.Time) {
				vo.Available = false
			}
		}
		items = append(items, vo)
	}
	return &pb.UserCouponPageReply{Total: total, List: items}, nil
}
