package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailPageQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailPageQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPageQueryLogic {
	return &OrderDetailPageQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailPageQueryLogic) OrderDetailPageQuery(in *pb.OrderDetailPageRequest) (*pb.OrderDetailPageReply, error) {
	// mobile 过滤因无 user 服务接入暂不支持；no_no 作为订单号精确过滤
	f := model.OrderDetailPageFilter{
		Id:           in.Id,
		OrderId:      in.NoNo,
		Status:       in.Status,
		RefundStatus: in.RefundStatus,
		PayChannel:   in.PayChannel,
		StartTime:    in.OrderStartTime,
		EndTime:      in.OrderEndTime,
		PageNo:       in.PageNo,
		PageSize:     in.PageSize,
		IsAsc:        in.IsAsc,
		SortBy:       in.SortBy,
	}

	details, total, err := l.svcCtx.OrderDetailModel.PageQuery(l.ctx, f)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细分页失败")
	}

	list := make([]*pb.OrderDetailPageVO, 0, len(details))
	for _, d := range details {
		refundStatus := nullInt64Value(d.RefundStatus)
		list = append(list, &pb.OrderDetailPageVO{
			Id:               d.Id,
			OrderId:          d.OrderId,
			CourseId:         d.CourseId,
			CourseName:       d.Name,
			Mobile:           "",
			Price:            d.Price,
			RealPayAmount:    d.RealPayAmount,
			DiscountAmount:   d.DiscountAmount,
			Status:           int32(d.Status),
			StatusDesc:       detailStatusDesc(d.Status),
			RefundStatus:     int32(refundStatus),
			RefundStatusDesc: refundStatusDesc(refundStatus),
			PayChannel:       d.PayChannel,
			CreateTime:       formatTime(d.CreateTime),
			FinishTime:       formatNullTime(d.CourseExpireTime),
		})
	}

	return &pb.OrderDetailPageReply{
		Total: total,
		Pages: calcPages(total, in.PageSize),
		List:  list,
	}, nil
}
