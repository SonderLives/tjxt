package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.EnrollCourseReply{}, nil
}
