package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailEnrollNumLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailEnrollNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailEnrollNumLogic {
	return &OrderDetailEnrollNumLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailEnrollNumLogic) OrderDetailEnrollNum(in *pb.EnrollNumRequest) (*pb.EnrollNumReply, error) {
	// course_id -> 报名人数
	items, err := l.svcCtx.OrderDetailModel.CountPaidByCourseIds(l.ctx, in.CourseIdList)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计课程报名人数失败")
	}
	if items == nil {
		items = make(map[int64]int64)
	}

	return &pb.EnrollNumReply{Items: items}, nil
}
