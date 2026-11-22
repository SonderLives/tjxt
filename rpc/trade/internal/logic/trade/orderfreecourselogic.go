package tradelogic

import (
	"context"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/internal/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"
	"tjxt/pkg/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderFreeCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderFreeCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderFreeCourseLogic {
	return &OrderFreeCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderFreeCourseLogic) OrderFreeCourse(in *pb.FreeCourseRequest) (*pb.PlaceOrderResultVO, error) {
	userId, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.CourseId <= 0 {
		return nil, errorx.BadRequestf("课程ID不能为空")
	}

	courseMap := fetchCourseMap(l.ctx, l.svcCtx, []int64{in.CourseId})

	// 免费课程金额为 0，下单即视为已支付
	order := &model.Order{
		Id:             idgen.NextID(),
		UserId:         userId,
		TotalAmount:    0,
		RealAmount:     0,
		DiscountAmount: 0,
		Status:         OrderStatusPaid,
		Creater:        userId,
		Updater:        userId,
		CreateTime:     now(),
	}
	if _, err = l.svcCtx.OrderModel.Insert(l.ctx, order); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "创建订单失败")
	}

	detail := &model.OrderDetail{
		Id:            idgen.NextID(),
		OrderId:       order.Id,
		UserId:        userId,
		CourseId:      in.CourseId,
		Name:          courseName(courseMap, in.CourseId),
		CoverUrl:      courseCover(courseMap, in.CourseId),
		Price:         0,
		RealPayAmount: 0,
		Status:        DetailStatusPaid,
		Creater:       userId,
		Updater:       userId,
		CreateTime:    now(),
	}
	if _, err = l.svcCtx.OrderDetailModel.Insert(l.ctx, detail); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "创建订单明细失败")
	}

	return &pb.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  0,
		Status:     int32(OrderStatusPaid),
		PayOutTime: now().Add(OrderPayTimeout).Unix(),
	}, nil
}
