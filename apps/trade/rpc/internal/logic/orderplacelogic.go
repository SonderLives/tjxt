package logic

import (
	"context"
	"time"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if len(in.CourseIds) == 0 {
		return nil, xerr.BadRequestf("课程ID不能为空")
	}

	courseMap := fetchCourseMap(l.ctx, l.svcCtx, in.CourseIds)

	var total int64
	for _, id := range in.CourseIds {
		total += coursePrice(courseMap, id)
	}

	// 优惠券暂未接入，实付金额等于课程价格之和
	order := &model.Order{
		Id:             nextID(),
		UserId:         userId,
		TotalAmount:    total,
		RealAmount:     total,
		DiscountAmount: 0,
		Status:         OrderStatusPending,
		Creater:        userId,
		Updater:        userId,
		CreateTime:     now(),
	}
	if _, err = l.svcCtx.OrderModel.Insert(l.ctx, order); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建订单失败")
	}

	for _, id := range in.CourseIds {
		price := coursePrice(courseMap, id)
		detail := &model.OrderDetail{
			Id:            nextID(),
			OrderId:       order.Id,
			UserId:        userId,
			CourseId:      id,
			Name:          courseName(courseMap, id),
			CoverUrl:      courseCover(courseMap, id),
			Price:         price,
			RealPayAmount: price,
			Status:        DetailStatusPending,
			Creater:       userId,
			Updater:       userId,
			CreateTime:    now(),
		}
		if _, err = l.svcCtx.OrderDetailModel.Insert(l.ctx, detail); err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "创建订单明细失败")
		}
	}

	return &pb.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  total,
		Status:     int32(OrderStatusPending),
		PayOutTime: now().Add(15 * time.Minute).Unix(),
	}, nil
}
