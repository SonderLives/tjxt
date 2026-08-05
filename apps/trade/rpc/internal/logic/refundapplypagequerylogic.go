package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyPageQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyPageQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyPageQueryLogic {
	return &RefundApplyPageQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyPageQueryLogic) RefundApplyPageQuery(in *pb.RefundApplyPageRequest) (*pb.RefundApplyPageReply, error) {
	pageNo, pageSize := in.PageNo, in.PageSize
	if pageNo < 1 {
		pageNo = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	// mobile 过滤因无 user 服务接入，暂不支持
	f := model.RefundApplyPageFilter{
		Id:            in.Id,
		OrderDetailId: in.OrderDetailId,
		OrderId:       in.OrderId,
		RefundStatus:  in.RefundStatus,
		StartTime:     in.ApplyStartTime,
		EndTime:       in.ApplyEndTime,
		PageNo:        pageNo,
		PageSize:      pageSize,
		IsAsc:         in.IsAsc,
		SortBy:        in.SortBy,
	}

	records, total, err := l.svcCtx.RefundApplyModel.PageQuery(l.ctx, f)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "分页查询退款申请失败")
	}

	list := make([]*pb.RefundApplyPageVO, 0, len(records))
	for _, ra := range records {
		list = append(list, toRefundApplyPageVO(ra))
	}

	return &pb.RefundApplyPageReply{
		Total: total,
		Pages: calcPages(total, pageSize),
		List:  list,
	}, nil
}
