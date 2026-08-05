package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SendNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendNoticeLogic {
	return &SendNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户收件箱
func (l *SendNoticeLogic) SendNotice(in *pb.SendNoticeReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
