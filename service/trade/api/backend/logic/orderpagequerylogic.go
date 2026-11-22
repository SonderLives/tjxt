package tradelogic

import (
	"context"

	"tjxt/api/trade/backend/svc"
	"tjxt/api/trade/backend/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"
	"tjxt/pkg/page"
	"tjxt/pkg/timex"

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
		return nil, err
	}

	orders, total, err := l.svcCtx.OrderModel.PageQueryByUser(l.ctx, userId, in.PageNo, in.PageSize, in.Status, in.NoNo, in.SortBy, in.IsAsc)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单分页失败")
	}

	list := make([]*pb.OrderPageVO, 0, len(orders))
	for _, order := range orders {
		vo := &pb.OrderPageVO{
			Id:          order.Id,
			Status:      order.Status,
			StatusDesc:  orderStatusDesc(order.Status),
			TotalAmount: order.TotalAmount,
			RealAmount:  order.RealAmount,
			CreateTime:  timex.Format(order.CreateTime),
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
		Pages: page.CalcPages(total, in.PageSize),
		List:  list,
	}, nil
}
