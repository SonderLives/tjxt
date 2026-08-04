package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"common/idgen"
	"common/xerr"
	"trade/internal/model"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// OrderService 订单业务接口
type OrderService interface {
	// PrePlaceOrder 预下单：创建待支付订单并返回确认信息。
	PrePlaceOrder(ctx context.Context, userId int64, courseIds []int64) (*types.OrderConfirmVO, error)
	// PlaceOrder 提交订单：校验并锁定金额。
	PlaceOrder(ctx context.Context, userId int64, req *types.PlaceOrderReq) (*types.PlaceOrderResultVO, error)
	// FreeCourse 免费课报名：直接生成已报名订单并发布开课事件。
	FreeCourse(ctx context.Context, userId, courseId int64) (*types.PlaceOrderResultVO, error)
	// CancelOrder 取消待支付订单。
	CancelOrder(ctx context.Context, userId, orderId int64) error
	// DeleteOrder 逻辑删除订单。
	DeleteOrder(ctx context.Context, userId, orderId int64) error
	// GetOrder 查询订单详情。
	GetOrder(ctx context.Context, userId, orderId int64) (*types.OrderVO, error)
	// PageOrders 分页查询用户订单。
	PageOrders(ctx context.Context, userId int64, pageNo, pageSize int64, status int64, isAsc bool) (*types.Page, error)
	// GetOrderStatus 查询订单支付状态。
	GetOrderStatus(ctx context.Context, userId, orderId int64) (*types.PlaceOrderResultVO, error)
}

// EventPublisher 订单事件发布抽象（支付/退款）。
type EventPublisher interface {
	PublishPay(ctx context.Context, orderID, userID int64, courseIDs []int64, finishTime time.Time) error
	PublishRefund(ctx context.Context, orderID, userID int64, courseIDs []int64, finishTime time.Time) error
}

type orderService struct {
	conn           sqlx.SqlConn
	orderModel     *model.OrderModel
	detailModel    *model.OrderDetailModel
	courseClient   CourseClient
	eventPublisher EventPublisher
}

// NewOrderService 创建订单业务服务。
func NewOrderService(conn sqlx.SqlConn, orderModel *model.OrderModel, detailModel *model.OrderDetailModel, courseClient CourseClient, publisher EventPublisher) OrderService {
	return &orderService{
		conn:           conn,
		orderModel:     orderModel,
		detailModel:    detailModel,
		courseClient:   courseClient,
		eventPublisher: publisher,
	}
}

// PrePlaceOrder 预下单。
func (s *orderService) PrePlaceOrder(ctx context.Context, userId int64, courseIds []int64) (*types.OrderConfirmVO, error) {
	if len(courseIds) == 0 {
		return nil, xerr.BadRequestf("请选择要购买的课程")
	}

	infos, err := s.courseClient.GetSimpleInfos(ctx, courseIds)
	if err != nil {
		return nil, err
	}
	if len(infos) != len(courseIds) {
		return nil, xerr.BadRequestf("部分课程不存在或不可购买")
	}

	now := time.Now()
	orderId := idgen.NextID()
	order := &model.Order{
		Id:             orderId,
		UserId:         userId,
		Status:         model.OrderStatusPendingPay,
		Message:        "待支付",
		TotalAmount:    0,
		RealAmount:     0,
		DiscountAmount: 0,
		PayChannel:     "",
		CreateTime:     now,
		Creater:        userId,
		Updater:        userId,
	}
	details := make([]*model.OrderDetail, 0, len(courseIds))
	courses := make([]types.OrderCourseVO, 0, len(courseIds))
	var total int64
	for _, cid := range courseIds {
		info := infos[cid]
		if info.Free {
			continue // 免费课不进入待支付订单
		}
		total += info.Price
		courses = append(courses, types.OrderCourseVO{
			Id:       info.Id,
			Name:     info.Name,
			CoverUrl: info.CoverUrl,
			Price:    info.Price,
		})
		details = append(details, &model.OrderDetail{
			Id:            idgen.NextID(),
			OrderId:       orderId,
			UserId:        userId,
			CourseId:      info.Id,
			Price:         info.Price,
			Name:          info.Name,
			CoverUrl:      info.CoverUrl,
			ValidDuration: sql.NullInt64{Int64: info.ValidDuration, Valid: info.ValidDuration > 0},
			Status:        model.OrderDetailStatusPendingPay,
			PayChannel:    "",
			Creater:       userId,
			Updater:       userId,
		})
	}
	if len(details) == 0 {
		return nil, xerr.BadRequestf("所选课程均为免费课程，可直接报名")
	}

	order.TotalAmount = total
	order.RealAmount = total
	err = s.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := s.orderModel.InsertTx(ctx, session, order); err != nil {
			return err
		}
		return s.detailModel.InsertBatchTx(ctx, session, details)
	})
	if err != nil {
		logx.Errorf("pre place order failed, user=%d err=%v", userId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "下单失败")
	}

	return &types.OrderConfirmVO{
		OrderId:     orderId,
		Courses:     courses,
		Discounts:   []types.CouponDiscountDTO{},
		TotalAmount: total,
	}, nil
}

// PlaceOrder 提交订单。
func (s *orderService) PlaceOrder(ctx context.Context, userId int64, req *types.PlaceOrderReq) (*types.PlaceOrderResultVO, error) {
	if req.OrderId == 0 {
		return nil, xerr.BadRequestf("订单 id 不能为空")
	}

	order, err := s.orderModel.FindById(ctx, req.OrderId)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("订单不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}
	if order.UserId != userId {
		return nil, xerr.Forbidden("无权操作该订单")
	}
	if order.Status != model.OrderStatusPendingPay {
		return nil, xerr.Conflict("订单状态不允许支付")
	}

	details, err := s.detailModel.ListByOrderId(ctx, order.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}

	// 允许用户在确认页调整课程：以提交的 courseIds 为准
	var real int64
	if len(req.CourseIds) > 0 {
		keep := make(map[int64]bool, len(req.CourseIds))
		for _, cid := range req.CourseIds {
			keep[cid] = true
		}
		filtered := details[:0]
		for _, d := range details {
			if keep[d.CourseId] {
				filtered = append(filtered, d)
				real += d.Price
			}
		}
		details = filtered
	} else {
		for i := range details {
			real += details[i].Price
		}
	}
	if len(details) == 0 {
		return nil, xerr.BadRequestf("订单中没有待购买课程")
	}

	// 优惠券暂未接入（promotion 服务未实现），抵扣为 0
	discount := int64(0)

	err = s.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := s.orderModel.UpdateAmount(ctx, order.Id, real, discount); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("place order failed, order=%d err=%v", order.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "提交订单失败")
	}

	payOutTime := time.Now().Add(15 * time.Minute)
	return &types.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  real,
		PayOutTime: payOutTime.Format(time.RFC3339),
		Status:     model.OrderStatusPendingPay,
	}, nil
}

// FreeCourse 免费课报名。
func (s *orderService) FreeCourse(ctx context.Context, userId, courseId int64) (*types.PlaceOrderResultVO, error) {
	if courseId == 0 {
		return nil, xerr.BadRequestf("课程 id 不能为空")
	}
	infos, err := s.courseClient.GetSimpleInfos(ctx, []int64{courseId})
	if err != nil {
		return nil, err
	}
	info, ok := infos[courseId]
	if !ok || info == nil {
		return nil, xerr.NotFound("课程不存在")
	}
	if !info.Free {
		return nil, xerr.BadRequestf("该课程不是免费课程")
	}

	// 幂等：已报名直接返回
	existing, err := s.detailModel.FindPaidByUserCourse(ctx, userId, courseId)
	if err == nil && existing != nil {
		return &types.PlaceOrderResultVO{
			OrderId:    existing.OrderId,
			PayAmount:  0,
			PayOutTime: existing.CreateTime.Format(time.RFC3339),
			Status:     model.OrderStatusEnrolled,
		}, nil
	}
	if err != nil && err != sql.ErrNoRows {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询报名信息失败")
	}

	now := time.Now()
	orderId := idgen.NextID()
	order := &model.Order{
		Id:             orderId,
		UserId:         userId,
		Status:         model.OrderStatusEnrolled,
		Message:        "免费报名",
		TotalAmount:    0,
		RealAmount:     0,
		DiscountAmount: 0,
		FinishTime:     sql.NullTime{Time: now, Valid: true},
		CreateTime:     now,
		Creater:        userId,
		Updater:        userId,
	}
	detail := &model.OrderDetail{
		Id:            idgen.NextID(),
		OrderId:       orderId,
		UserId:        userId,
		CourseId:      info.Id,
		Price:         0,
		Name:          info.Name,
		CoverUrl:      info.CoverUrl,
		ValidDuration: sql.NullInt64{Int64: info.ValidDuration, Valid: info.ValidDuration > 0},
		Status:        model.OrderDetailStatusEnrolled,
		Creater:       userId,
		Updater:       userId,
	}

	err = s.conn.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if err := s.orderModel.InsertTx(ctx, session, order); err != nil {
			return err
		}
		return s.detailModel.InsertTx(ctx, session, detail)
	})
	if err != nil {
		logx.Errorf("free course enroll failed, user=%d course=%d err=%v", userId, courseId, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "免费报名失败")
	}

	// 发布开课事件，通知 learning 服务为用户开通课程
	if err := s.eventPublisher.PublishPay(ctx, orderId, userId, []int64{courseId}, now); err != nil {
		logx.Errorf("publish free course pay event failed: %v", err)
	}

	return &types.PlaceOrderResultVO{
		OrderId:    orderId,
		PayAmount:  0,
		PayOutTime: now.Format(time.RFC3339),
		Status:     model.OrderStatusEnrolled,
	}, nil
}

// CancelOrder 取消待支付订单。
func (s *orderService) CancelOrder(ctx context.Context, userId, orderId int64) error {
	affected, err := s.orderModel.CancelPending(ctx, orderId, userId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "取消订单失败")
	}
	if affected == 0 {
		// 可能不存在或状态不允许，返回冲突
		return xerr.Conflict("订单不存在或当前状态不允许取消")
	}
	return nil
}

// DeleteOrder 逻辑删除订单。
func (s *orderService) DeleteOrder(ctx context.Context, userId, orderId int64) error {
	affected, err := s.orderModel.SoftDelete(ctx, orderId, userId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "删除订单失败")
	}
	if affected == 0 {
		return xerr.Conflict("订单不存在或当前状态不允许删除")
	}
	return nil
}

// GetOrder 订单详情。
func (s *orderService) GetOrder(ctx context.Context, userId, orderId int64) (*types.OrderVO, error) {
	order, err := s.orderModel.FindById(ctx, orderId)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("订单不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}
	if order.UserId != userId {
		return nil, xerr.Forbidden("无权查看该订单")
	}

	details, err := s.detailModel.ListByOrderId(ctx, order.Id)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}
	detailVOs := make([]types.OrderDetailVO, 0, len(details))
	for i := range details {
		detailVOs = append(detailVOs, detailToVO(&details[i]))
	}

	return &types.OrderVO{
		Id:             order.Id,
		Status:         order.Status,
		StatusDesc:     orderStatusDesc(order.Status),
		Message:        order.Message,
		TotalAmount:    order.TotalAmount,
		RealAmount:     order.RealAmount,
		DiscountAmount: order.DiscountAmount,
		CouponDesc:     parseCouponDesc(order.CouponIds),
		CreateTime:     order.CreateTime.Format(time.RFC3339),
		Details:        detailVOs,
		ProgressNodes:  orderProgressNodes(order),
	}, nil
}

// PageOrders 分页查询用户订单。
func (s *orderService) PageOrders(ctx context.Context, userId int64, pageNo, pageSize int64, status int64, isAsc bool) (*types.Page, error) {
	offset, limit := normalizePage(pageNo, pageSize)
	orders, total, err := s.orderModel.ListByUser(ctx, userId, status, offset, limit, isAsc)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}

	orderIds := make([]int64, 0, len(orders))
	for i := range orders {
		orderIds = append(orderIds, orders[i].Id)
	}
	detailsByOrder := make(map[int64][]types.OrderDetailVO)
	if len(orderIds) > 0 {
		details, err := s.detailModel.ListByOrderIds(ctx, orderIds)
		if err != nil {
			return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
		}
		for i := range details {
			oid := details[i].OrderId
			detailsByOrder[oid] = append(detailsByOrder[oid], detailToVO(&details[i]))
		}
	}

	list := make([]types.OrderPageVO, 0, len(orders))
	for i := range orders {
		o := &orders[i]
		list = append(list, types.OrderPageVO{
			Id:          o.Id,
			Status:      o.Status,
			StatusDesc:  orderStatusDesc(o.Status),
			TotalAmount: o.TotalAmount,
			RealAmount:  o.RealAmount,
			CreateTime:  o.CreateTime.Format(time.RFC3339),
			Details:     detailsByOrder[o.Id],
		})
	}
	return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
}

// GetOrderStatus 查询订单支付状态。
func (s *orderService) GetOrderStatus(ctx context.Context, userId, orderId int64) (*types.PlaceOrderResultVO, error) {
	order, err := s.orderModel.FindById(ctx, orderId)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("订单不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}
	if order.UserId != userId {
		return nil, xerr.Forbidden("无权查看该订单")
	}

	payOutTime := order.CloseTime
	if !payOutTime.Valid && order.Status == model.OrderStatusPendingPay {
		payOutTime = sql.NullTime{Time: order.CreateTime.Add(15 * time.Minute), Valid: true}
	}
	return &types.PlaceOrderResultVO{
		OrderId:    order.Id,
		PayAmount:  order.RealAmount,
		PayOutTime: formatNullTime(payOutTime),
		Status:     order.Status,
	}, nil
}

func normalizePage(pageNo, pageSize int64) (offset, limit int64) {
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return (pageNo - 1) * pageSize, pageSize
}

func calcPages(total, pageSize int64) int64 {
	if pageSize <= 0 || total <= 0 {
		return 0
	}
	pages := total / pageSize
	if total%pageSize != 0 {
		pages++
	}
	return pages
}

// parseCouponDesc 解析订单优惠券描述。
func parseCouponDesc(couponIds sql.NullString) string {
	if !couponIds.Valid || couponIds.String == "" {
		return ""
	}
	var ids []int64
	if err := json.Unmarshal([]byte(couponIds.String), &ids); err != nil || len(ids) == 0 {
		return ""
	}
	return "优惠券抵扣"
}
