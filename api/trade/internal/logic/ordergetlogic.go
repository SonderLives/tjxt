// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderGetLogic {
	return &OrderGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderGetLogic) OrderGet(req *types.OrderIdReq) (resp *types.OrderVO, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	resp = &types.OrderVO{
		Id:             reply.Id,
		CreateTime:     reply.CreateTime,
		TotalAmount:    reply.TotalAmount,
		RealAmount:     reply.RealAmount,
		DiscountAmount: reply.DiscountAmount,
		Status:         int64(reply.Status),
		StatusDesc:     reply.StatusDesc,
		Message:        reply.Message,
		CouponDesc:     reply.CouponDesc,
		Details:        make([]types.OrderDetailItemVO, 0, len(reply.Details)),
		ProgressNodes:  make([]types.OrderProgressNodeVO, 0, len(reply.ProgressNodes)),
	}
	for _, d := range reply.Details {
		resp.Details = append(resp.Details, types.OrderDetailItemVO{
			Id:            d.Id,
			CourseId:      d.CourseId,
			CourseName:    d.CourseName,
			CoverUrl:      d.CoverUrl,
			Price:         d.Price,
			RealPayAmount: d.RealPayAmount,
			Status:        int64(d.Status),
			RefundStatus:  int64(d.RefundStatus),
			CanRefund:     d.CanRefund,
		})
	}
	for _, n := range reply.ProgressNodes {
		resp.ProgressNodes = append(resp.ProgressNodes, types.OrderProgressNodeVO{
			Name:   n.Name,
			Desc:   n.Desc,
			Status: int64(n.Status),
			Time:   n.Time,
		})
	}
	return resp, nil
}
