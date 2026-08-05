package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/utils/page"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponPageLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPageLogic {
	return &CouponPageLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponPage 管理端分页查询优惠券，支持名称模糊、状态与类型过滤。
func (l *CouponPageLogic) CouponPage(in *pb.CouponPageRequest) (*pb.CouponPageReply, error) {
	var pageNo, pageSize int64
	if in.Page != nil {
		pageNo, pageSize = in.Page.PageNo, in.Page.PageSize
	}
	offset, limit := page.Normalize(pageNo, pageSize)

	list, total, err := l.svcCtx.CouponModel.FindPage(l.ctx, in.Name, in.Status, in.Type, offset, limit)
	if err != nil {
		l.Errorf("page coupon failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券失败")
	}

	items := make([]*pb.CouponPageVO, 0, len(list))
	for _, c := range list {
		items = append(items, toCouponPageVO(c))
	}
	return &pb.CouponPageReply{Total: total, List: items}, nil
}
