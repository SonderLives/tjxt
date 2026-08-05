package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPublicNoticesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPublicNoticesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPublicNoticesLogic {
	return &ListPublicNoticesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPublicNoticesLogic) ListPublicNotices(in *pb.PageReq) (*pb.PublicNoticeListReply, error) {
	// todo: add your logic here and delete this line

	return &pb.PublicNoticeListReply{}, nil
}
