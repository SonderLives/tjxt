// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package svc

import (
	"tjxt/apps/trade/api/internal/config"
	"tjxt/apps/trade/api/internal/model"
	"tjxt/apps/trade/api/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	// 数据访问
	CartModel       *model.CartModel
	OrderModel      *model.OrderModel
	DetailModel     *model.OrderDetailModel
	RefundModel     *model.RefundApplyModel
	PayChannelModel *model.PayChannelModel

	// 跨服务客户端
	CourseClient service.CourseClient
	UserClient   service.UserClient

	// 事件发布
	EventPublisher service.EventPublisher

	// 业务服务
	CartService        service.CartService
	OrderService       service.OrderService
	PayService         service.PayService
	RefundService      service.RefundService
	OrderDetailService service.OrderDetailService
}

func NewServiceContext(c config.Config, publisher service.EventPublisher) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	// 数据访问层
	cartModel := model.NewCartModel(conn)
	orderModel := model.NewOrderModel(conn)
	detailModel := model.NewOrderDetailModel(conn)
	refundModel := model.NewRefundApplyModel(conn)
	payChannelModel := model.NewPayChannelModel(conn)

	// 跨服务客户端（无可用时返回空实现，避免 panic，接口返回依赖不可用）
	courseClient := service.NewCourseHTTPClient(c.CourseService.BaseURL, c.CourseService.Timeout)
	userClient := service.NewUserHTTPClient(c.UserService.BaseURL, c.UserService.Timeout)

	// 事件发布器（producer 在 main 中构建后注入）
	eventPublisher := service.NewMQEventPublisher(nil, c.Exchange, c.PayRoutingKey, c.RefundRoutingKey)
	if publisher != nil {
		eventPublisher = publisher
	}

	sc := &ServiceContext{
		Config:             c,
		CartModel:          cartModel,
		OrderModel:         orderModel,
		DetailModel:        detailModel,
		RefundModel:        refundModel,
		PayChannelModel:    payChannelModel,
		CourseClient:       courseClient,
		UserClient:         userClient,
		EventPublisher:     eventPublisher,
		CartService:        service.NewCartService(cartModel, courseClient),
		OrderService:       service.NewOrderService(conn, orderModel, detailModel, courseClient, eventPublisher),
		PayService:         service.NewPayService(orderModel, detailModel, payChannelModel, eventPublisher),
		RefundService:      service.NewRefundService(refundModel, orderModel, detailModel, eventPublisher, userClient),
		OrderDetailService: service.NewOrderDetailService(detailModel, orderModel, refundModel, userClient),
	}

	return sc
}
