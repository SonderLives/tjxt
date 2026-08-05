package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderPageQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderPageQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPageQueryLogic {
	return &OrderPageQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderPageQueryLogic) OrderPageQuery(in *pb.OrderPageRequest) (*pb.OrderPageReply, error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}

	orders, total, err := l.svcCtx.OrderModel.PageQueryByUser(l.ctx, userId, in.PageNo, in.PageSize, in.Status, in.NoNo, in.SortBy, in.IsAsc)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单分页失败")
	}

	list := make([]*pb.OrderPageVO, 0, len(orders))
	for _, order := range orders {
		vo := &pb.OrderPageVO{
			Id:          order.Id,
			Status:      order.Status,
			StatusDesc:  orderStatusDesc(order.Status),
			TotalAmount: order.TotalAmount,
			RealAmount:  order.RealAmount,
			CreateTime:  formatTime(order.CreateTime),
		}
		details, e := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, order.Id)
		if e == nil {
			for _, d := range details {
				vo.Details = append(vo.Details, toOrderDetailItemVO(d))
			}
		}
		list = append(list, vo)
	}

	return &pb.OrderPageReply{
		Total: total,
		Pages: calcPages(total, in.PageSize),
		List:  list,
	}, nil
}
