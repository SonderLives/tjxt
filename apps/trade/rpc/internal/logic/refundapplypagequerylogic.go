package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"

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
	// todo: add your logic here and delete this line

	return &pb.RefundApplyPageReply{}, nil
}
