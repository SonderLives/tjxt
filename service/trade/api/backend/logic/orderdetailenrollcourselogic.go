package tradelogic

import (
	"context"

	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/errorx"

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
		return nil, errorx.Wrap(err, errorx.CodeInternal, "统计学员报名课程数失败")
	}
	if items == nil {
		items = make(map[int64]int64)
	}

	return &pb.EnrollCourseReply{Items: items}, nil
}
