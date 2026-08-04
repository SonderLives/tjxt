package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderFreeCourseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderFreeCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderFreeCourseLogic {
	return &OrderFreeCourseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderFreeCourseLogic) OrderFreeCourse(in *pb.FreeCourseRequest) (*pb.PlaceOrderResultVO, error) {
	// todo: add your logic here and delete this line

	return &pb.PlaceOrderResultVO{}, nil
}
