package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailEnrollCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailEnrollCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailEnrollCourseLogic {
	return &OrderDetailEnrollCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailEnrollCourseLogic) OrderDetailEnrollCourse(in *pb.EnrollCourseRequest) (*pb.EnrollCourseReply, error) {
	// student_id -> 已报名课程数
	items, err := l.svcCtx.OrderDetailModel.CountPaidByUserIds(l.ctx, in.StudentIds)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "统计学员报名课程数失败")
	}
	if items == nil {
		items = make(map[int64]int64)
	}

	return &pb.EnrollCourseReply{Items: items}, nil
}
