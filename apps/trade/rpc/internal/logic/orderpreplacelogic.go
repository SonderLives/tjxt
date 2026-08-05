package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPrePlaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderPrePlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPrePlaceLogic {
	return &OrderPrePlaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 订单 =====
func (l *OrderPrePlaceLogic) OrderPrePlace(in *pb.PrePlaceRequest) (*pb.OrderConfirmVO, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if len(in.CourseIds) == 0 {
		return nil, xerr.BadRequestf("课程ID不能为空")
	}

	courseMap := fetchCourseMap(l.ctx, l.svcCtx, in.CourseIds)

	var total int64
	courses := make([]*pb.OrderCourseVO, 0, len(in.CourseIds))
	for _, id := range in.CourseIds {
		price := coursePrice(courseMap, id)
		total += price
		courses = append(courses, &pb.OrderCourseVO{
			Id:       id,
			Name:     courseName(courseMap, id),
			CoverUrl: courseCover(courseMap, id),
			Price:    price,
		})
	}

	// 优惠券无服务接入，折扣列表返回空
	return &pb.OrderConfirmVO{
		OrderId:     0,
		TotalAmount: total,
		Courses:     courses,
		Discounts:   []*pb.CouponDiscountDTO{},
	}, nil
}
