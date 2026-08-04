package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailCourseCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailCourseCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailCourseCheckLogic {
	return &OrderDetailCourseCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailCourseCheckLogic) OrderDetailCourseCheck(in *pb.IdRequest) (*pb.BoolReply, error) {
	// todo: add your logic here and delete this line

	return &pb.BoolReply{}, nil
}
