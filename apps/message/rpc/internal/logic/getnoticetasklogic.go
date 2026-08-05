package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetNoticeTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetNoticeTaskLogic {
	return &GetNoticeTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetNoticeTaskLogic) GetNoticeTask(in *pb.IdReq) (*pb.NoticeTaskVO, error) {
	// todo: add your logic here and delete this line

	return &pb.NoticeTaskVO{}, nil
}
