package promotionlogic

import (
	"context"

	"tjxt/app/promotion/rpc/internal/svc"
	"tjxt/app/promotion/rpc/pb"
	"tjxt/pkg/page"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponCodePageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponCodePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponCodePageLogic {
	return &CouponCodePageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponCodePage 管理端分页查询兑换码。
func (l *CouponCodePageLogic) CouponCodePage(in *pb.CouponCodePageRequest) (*pb.CouponCodePageReply, error) {
	var pageNo, pageSize int64
	if in.Page != nil {
		pageNo, pageSize = in.Page.PageNo, in.Page.PageSize
	}
	offset, limit := page.Normalize(pageNo, pageSize)

	list, total, err := l.svcCtx.CouponCodeModel.FindPageByCoupon(l.ctx, in.CouponId, in.Status, offset, limit)
	if err != nil {
		l.Errorf("page coupon code failed, couponId=%d, err=%v", in.CouponId, err)
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询兑换码失败")
	}

	items := make([]*pb.ExchangeCodeVO, 0, len(list))
	for _, cc := range list {
		items = append(items, &pb.ExchangeCodeVO{Id: cc.Id, Code: cc.Code})
	}
	return &pb.CouponCodePageReply{Total: total, List: items}, nil
}
