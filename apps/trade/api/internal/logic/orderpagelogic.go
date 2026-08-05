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

type OrderPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderPageLogic {
	return &OrderPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderPageLogic) OrderPage(req *types.OrderPageReq) (resp *types.OrderPageReply, err error) {
	reply, err := l.svcCtx.TradeRpc.OrderPageQuery(l.ctx, &pb.OrderPageRequest{
		PageNo:   req.PageNo,
		PageSize: req.PageSize,
		IsAsc:    req.IsAsc,
		SortBy:   req.SortBy,
		NoNo:     req.NoNo,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}

	resp = &types.OrderPageReply{
		Total: reply.Total,
		Pages: reply.Pages,
		List:  make([]types.OrderPageVO, 0, len(reply.List)),
	}
	for _, o := range reply.List {
		vo := types.OrderPageVO{
			Id:          o.Id,
			Status:      o.Status,
			StatusDesc:  o.StatusDesc,
			TotalAmount: o.TotalAmount,
			RealAmount:  o.RealAmount,
			CreateTime:  o.CreateTime,
			Details:     make([]types.OrderDetailItemVO, 0, len(o.Details)),
		}
		for _, d := range o.Details {
			vo.Details = append(vo.Details, types.OrderDetailItemVO{
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
		resp.List = append(resp.List, vo)
	}
	return resp, nil
}
