// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strconv"
	"strings"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPrePlaceLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPrePlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPrePlaceLogic {
	return &OrderPrePlaceLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderPrePlaceLogic) OrderPrePlace(req *types.PrePlaceOrderReq) (resp *types.OrderConfirmVO, err error) {
	var ids []int64
	for _, s := range strings.Split(req.CourseIds, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, parseErr := strconv.ParseInt(s, 10, 64)
		if parseErr != nil {
			continue
		}
		ids = append(ids, id)
	}

	reply, err := l.svcCtx.TradeRpc.OrderPrePlace(l.ctx, &pb.PrePlaceRequest{CourseIds: ids})
	if err != nil {
		return nil, err
	}

	resp = &types.OrderConfirmVO{
		OrderId:     reply.OrderId,
		TotalAmount: reply.TotalAmount,
		Courses:     make([]types.OrderCourseVO, 0, len(reply.Courses)),
		Discounts:   make([]types.CouponDiscountVO, 0, len(reply.Discounts)),
	}
	for _, c := range reply.Courses {
		resp.Courses = append(resp.Courses, types.OrderCourseVO{
			Id:       c.Id,
			Name:     c.Name,
			CoverUrl: c.CoverUrl,
			Price:    c.Price,
		})
	}
	for _, d := range reply.Discounts {
		resp.Discounts = append(resp.Discounts, types.CouponDiscountVO{
			Id:             d.Id,
			Name:           d.Name,
			DiscountAmount: d.DiscountAmount,
			RuleDesc:       d.RuleDesc,
		})
	}
	return resp, nil
}
