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

type OrderDetailPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPageLogic {
	return &OrderDetailPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailPageLogic) OrderDetailPage(req *types.OrderDetailPageReq) (resp *types.OrderDetailPageReply, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderDetailPageQuery(l.ctx, &pb.OrderDetailPageRequest{
		PageNo:         req.PageNo,
		PageSize:       req.PageSize,
		IsAsc:          req.IsAsc,
		SortBy:         req.SortBy,
		Id:             req.Id,
		Mobile:         req.Mobile,
		Status:         req.Status,
		RefundStatus:   req.RefundStatus,
		PayChannel:     req.PayChannel,
		OrderStartTime: req.OrderStartTime,
		OrderEndTime:   req.OrderEndTime,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.OrderDetailPageReply{
		Total: reply.Total,
		Pages: reply.Pages,
		List:  make([]types.OrderDetailPageVO, 0, len(reply.List)),
	}
	for _, d := range reply.List {
		resp.List = append(resp.List, types.OrderDetailPageVO{
			Id:               d.Id,
			OrderId:          d.OrderId,
			CourseId:         d.CourseId,
			CourseName:       d.CourseName,
			Mobile:           d.Mobile,
			Price:            d.Price,
			RealPayAmount:    d.RealPayAmount,
			DiscountAmount:   d.DiscountAmount,
			Status:           int64(d.Status),
			StatusDesc:       d.StatusDesc,
			RefundStatus:     int64(d.RefundStatus),
			RefundStatusDesc: d.RefundStatusDesc,
			PayChannel:       d.PayChannel,
			CreateTime:       d.CreateTime,
			FinishTime:       d.FinishTime,
		})
	}
	return resp, nil
}
