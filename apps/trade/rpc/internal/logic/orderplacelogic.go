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
	"github.com/zeromicro/go-zero/core/stores/sqlx"
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

	details := make([]*model.OrderDetail, 0, len(in.CourseIds))
	for _, id := range in.CourseIds {
		price := coursePrice(courseMap, id)
		details = append(details, &model.OrderDetail{
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
		})
	}

	// order + details 同事务写入：中途失败整体回滚，避免留下无明细的脏订单。
	// 新行尚无缓存条目，事务内直接 Exec 即可，无需失效缓存。
	if err = l.svcCtx.DB.TransactCtx(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx,
			"insert into `order` (id, user_id, total_amount, real_amount, discount_amount, status, creater, updater, create_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?)",
			order.Id, order.UserId, order.TotalAmount, order.RealAmount, order.DiscountAmount,
			order.Status, order.Creater, order.Updater, order.CreateTime); err != nil {
			return err
		}
		for _, d := range details {
			if _, err := session.ExecCtx(ctx,
				"insert into `order_detail` (id, order_id, user_id, course_id, name, cover_url, price, real_pay_amount, status, creater, updater, create_time) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
				d.Id, d.OrderId, d.UserId, d.CourseId, d.Name, d.CoverUrl,
				d.Price, d.RealPayAmount, d.Status, d.Creater, d.Updater, d.CreateTime); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建订单失败")
	}

	return &pb.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  total,
		Status:     int32(OrderStatusPending),
		PayOutTime: now().Add(15 * time.Minute).Unix(),
	}, nil
}
