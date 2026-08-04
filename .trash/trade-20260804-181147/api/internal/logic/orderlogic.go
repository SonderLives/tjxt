package logic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/response"
	"tjxt/pkg/xerr"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// canJoinFreeCourse 判断角色是否可报名免费课程（学员或助教）。
func canJoinFreeCourse(role string) bool {
	switch role {
	case auth.RoleStudent, auth.RoleTeacher, auth.RoleStaff:
		return true
	}
	return false
}

// ============ 预下单 ============

type PrePlaceOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPrePlaceOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PrePlaceOrderLogic {
	return &PrePlaceOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PrePlaceOrderLogic) PrePlaceOrder(req *types.PrePlaceOrderReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	courseIDs := parseIDs(req.CourseIds)
	if len(courseIDs) == 0 {
		return nil, errBadRequest("请选择课程")
	}
	confirm, err := l.svcCtx.OrderService.PrePlaceOrder(l.ctx, userID, courseIDs)
	if err != nil {
		return nil, err
	}
	return result.OkData(confirm), nil
}

// ============ 提交订单 ============

type PlaceOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPlaceOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PlaceOrderLogic {
	return &PlaceOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PlaceOrderLogic) PlaceOrder(req *types.PlaceOrderReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	resultVO, err := l.svcCtx.OrderService.PlaceOrder(l.ctx, userID, req)
	if err != nil {
		return nil, err
	}
	return result.OkData(resultVO), nil
}

// ============ 免费课报名 ============

type FreeCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFreeCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FreeCourseLogic {
	return &FreeCourseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *FreeCourseLogic) FreeCourse(req *types.FreeCourseReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if !canJoinFreeCourse(auth.RoleFromCtx(l.ctx)) {
		return nil, xerr.Forbidden("仅学员或助教可报名免费课程")
	}
	resultVO, err := l.svcCtx.OrderService.FreeCourse(l.ctx, userID, req.CourseId)
	if err != nil {
		return nil, err
	}
	return result.OkData(resultVO), nil
}

// ============ 订单分页 ============

type OrderPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPageLogic {
	return &OrderPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderPageLogic) OrderPage(req *types.OrderPageReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	page, err := l.svcCtx.OrderService.PageOrders(l.ctx, userID, req.PageNo, req.PageSize, req.Status, req.IsAsc)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 订单详情 ============

type OrderDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailLogic {
	return &OrderDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderDetailLogic) OrderDetail(req *types.OrderIdReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	order, err := l.svcCtx.OrderService.GetOrder(l.ctx, userID, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(order), nil
}

// ============ 删除订单 ============

type DeleteOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOrderLogic {
	return &DeleteOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *DeleteOrderLogic) DeleteOrder(req *types.OrderIdReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.OrderService.DeleteOrder(l.ctx, userID, req.Id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 取消订单 ============

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CancelOrderLogic) CancelOrder(req *types.OrderIdReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.OrderService.CancelOrder(l.ctx, userID, req.Id); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 订单支付状态 ============

type OrderStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderStatusLogic {
	return &OrderStatusLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderStatusLogic) OrderStatus(req *types.OrderIdReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	statusVO, err := l.svcCtx.OrderService.GetOrderStatus(l.ctx, userID, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(statusVO), nil
}
