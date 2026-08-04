package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.EnrollNumReply{}, nil
}
