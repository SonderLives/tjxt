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

type RefundApplyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyPageLogic {
	return &RefundApplyPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundApplyPageLogic) RefundApplyPage(req *types.RefundApplyPageReq) (resp *types.RefundApplyPageReply, err error) {
	reply, err := l.svcCtx.TradeRpc.RefundApplyPageQuery(l.ctx, &pb.RefundApplyPageRequest{
		PageNo:         req.PageNo,
		PageSize:       req.PageSize,
		IsAsc:          req.IsAsc,
		SortBy:         req.SortBy,
		Id:             req.Id,
		OrderDetailId:  req.OrderDetailId,
		OrderId:        req.OrderId,
		RefundStatus:   req.RefundStatus,
		Mobile:         req.Mobile,
		ApplyStartTime: req.ApplyStartTime,
		ApplyEndTime:   req.ApplyEndTime,
	})
	if err != nil {
		return nil, err
	}

	list := make([]types.RefundApplyPageVO, 0, len(reply.List))
	for _, item := range reply.List {
		list = append(list, types.RefundApplyPageVO{
			Id:            item.Id,
			OrderId:       item.OrderId,
			OrderDetailId: item.OrderDetailId,
			Price:         item.Price,
			RefundAmount:  item.RefundAmount,
			Status:        int64(item.Status),
			StatusDesc:    item.StatusDesc,
			Mobile:        item.Mobile,
			StudentName:   item.StudentName,
			RefundReason:  item.RefundReason,
			CreateTime:    item.CreateTime,
			ApproveTime:   item.ApproveTime,
		})
	}
	return &types.RefundApplyPageReply{
		Total: reply.Total,
		Pages: reply.Pages,
		List:  list,
	}, nil
}
