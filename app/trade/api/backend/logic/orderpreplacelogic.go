package tradelogic

import (
	"context"
	"strings"

	promotionclient "tjxt/service/promotion/rpc/client/promotion"
	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/auth"
	"tjxt/common/errorx"

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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if len(in.CourseIds) == 0 {
		return nil, errorx.BadRequestf("课程ID不能为空")
	}

	courseMap := fetchCourseMap(l.ctx, l.svcCtx, in.CourseIds)

	var total int64
	courses := make([]*pb.OrderCourseVO, 0, len(in.CourseIds))
	courseList := make([]*promotionclient.OrderCourseDTO, 0, len(in.CourseIds))
	for _, id := range in.CourseIds {
		price := coursePrice(courseMap, id)
		total += price
		courses = append(courses, &pb.OrderCourseVO{
			Id:       id,
			Name:     courseName(courseMap, id),
			CoverUrl: courseCover(courseMap, id),
			Price:    price,
		})
		courseList = append(courseList, &promotionclient.OrderCourseDTO{
			CateId: courseThirdCate(courseMap, id),
			Id:     id,
			Price:  price,
		})
	}

	// 优惠券可用方案：promotion 按课程范围/门槛计算最优组合；多券组合方案
	// 在展示层取首券 id（下单时传整组 coupon_ids 即可）
	discounts := make([]*pb.CouponDiscountDTO, 0)
	if reply, err := l.svcCtx.PromotionRpc.UserCouponAvailable(l.ctx, &promotionclient.OrderCourseListRequest{
		CourseList: courseList,
		UserId:     userId,
	}); err != nil {
		// 优惠方案查询失败不阻塞下单确认，仅告警并返回无折扣
		l.Errorf("query available coupons failed, userId=%d: %v", userId, err)
	} else {
		for _, d := range reply.List {
			firstID := int64(0)
			if len(d.Ids) > 0 {
				firstID = d.Ids[0]
			}
			discounts = append(discounts, &pb.CouponDiscountDTO{
				Id:             firstID,
				Name:           strings.Join(d.Rules, " + "),
				DiscountAmount: d.DiscountAmount,
				RuleDesc:       strings.Join(d.Rules, "; "),
			})
		}
	}

	return &pb.OrderConfirmVO{
		OrderId:     0,
		TotalAmount: total,
		Courses:     courses,
		Discounts:   discounts,
	}, nil
}
