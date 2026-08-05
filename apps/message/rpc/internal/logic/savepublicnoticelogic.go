package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePublicNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSavePublicNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePublicNoticeLogic {
	return &SavePublicNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公告
func (l *SavePublicNoticeLogic) SavePublicNotice(in *pb.PublicNoticeSaveReq) (*pb.IdReply, error) {
	// todo: add your logic here and delete this line

	return &pb.IdReply{}, nil
}
