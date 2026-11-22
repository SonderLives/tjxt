package tradelogic

import (
	"context"
	"database/sql"
	"encoding/json"

	promotionclient "tjxt/rpc/promotion/client/promotion"
	"tjxt/api/trade/backend/model"
	"tjxt/api/trade/backend/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"
	"tjxt/pkg/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPlaceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderPlaceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPlaceLogic {
	return &OrderPlaceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderPlaceLogic) OrderPlace(in *pb.PlaceOrderRequest) (*pb.PlaceOrderResultVO, error) {
	userId, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	if len(in.CourseIds) == 0 {
		return nil, errorx.BadRequestf("课程ID不能为空")
	}

	courseMap := fetchCourseMap(l.ctx, l.svcCtx, in.CourseIds)

	var total int64
	courseList := make([]*promotionclient.OrderCourseDTO, 0, len(in.CourseIds))
	for _, id := range in.CourseIds {
		total += coursePrice(courseMap, id)
		courseList = append(courseList, &promotionclient.OrderCourseDTO{
			CateId: courseThirdCate(courseMap, id),
			Id:     id,
			Price:  coursePrice(courseMap, id),
		})
	}

	// 优惠券折扣：选了券则调 promotion 计算总优惠与每课分摊
	var discountAmount int64
	discountDetail := map[int64]int64{}
	if len(in.CouponIds) > 0 {
		reply, derr := l.svcCtx.PromotionRpc.UserCouponDiscount(l.ctx, &promotionclient.OrderCouponDTO{
			CourseList:    courseList,
			UserCouponIds: in.CouponIds,
			UserId:        userId,
		})
		if derr != nil {
			return nil, errorx.Wrap(derr, errorx.CodeInternal, "计算优惠券折扣失败")
		}
		discountAmount = reply.DiscountAmount
		if discountAmount > total {
			discountAmount = total
		}
		for cid, amt := range reply.DiscountDetail {
			discountDetail[cid] = amt
		}
	}

	order := &model.Order{
		Id:             idgen.NextID(),
		UserId:         userId,
		TotalAmount:    total,
		RealAmount:     total - discountAmount,
		DiscountAmount: discountAmount,
		Status:         OrderStatusPending,
		Creater:        userId,
		Updater:        userId,
		CreateTime:     now(),
	}
	if len(in.CouponIds) > 0 {
		if b, jerr := json.Marshal(in.CouponIds); jerr == nil {
			order.CouponIds = sql.NullString{String: string(b), Valid: true}
		}
	}

	details := make([]*model.OrderDetail, 0, len(in.CourseIds))
	for _, id := range in.CourseIds {
		price := coursePrice(courseMap, id)
		courseDiscount := discountDetail[id]
		realPay := price - courseDiscount
		if realPay < 0 {
			realPay = 0
		}
		details = append(details, &model.OrderDetail{
			Id:             idgen.NextID(),
			OrderId:        order.Id,
			UserId:         userId,
			CourseId:       id,
			Name:           courseName(courseMap, id),
			CoverUrl:       courseCover(courseMap, id),
			Price:          price,
			RealPayAmount:  realPay,
			DiscountAmount: courseDiscount,
			Status:         DetailStatusPending,
			Creater:        userId,
			Updater:        userId,
			CreateTime:     now(),
		})
	}

	// order + details 同事务写入（事务下沉至 model，logic 只编排）：
	// 中途失败整体回滚，避免留下无明细的脏订单
	if err = l.svcCtx.OrderModel.PlaceOrderTx(l.ctx, order, details); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "创建订单失败")
	}

	// 核销优惠券：失败不影响订单创建（折扣已生效），券状态由过期/对账兜底
	if len(in.CouponIds) > 0 {
		if _, uerr := l.svcCtx.PromotionRpc.UserCouponUse(l.ctx, &promotionclient.IdsRequest{
			Ids:     in.CouponIds,
			UserId:  userId,
			OrderId: order.Id,
		}); uerr != nil {
			l.Errorf("use coupons failed, orderId=%d couponIds=%v: %v", order.Id, in.CouponIds, uerr)
		}
	}

	return &pb.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  order.RealAmount,
		Status:     int32(OrderStatusPending),
		PayOutTime: now().Add(OrderPayTimeout).Unix(),
	}, nil
}
