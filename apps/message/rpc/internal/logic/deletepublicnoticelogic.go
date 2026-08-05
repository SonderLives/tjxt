package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeletePublicNoticeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeletePublicNoticeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeletePublicNoticeLogic {
	return &DeletePublicNoticeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 公告 物理删除
func (l *DeletePublicNoticeLogic) DeletePublicNotice(in *pb.IdReq) (*pb.Empty, error) {
	if err := l.svcCtx.PublicNoticeModel.Delete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除公告失败")
	}
	return &pb.Empty{}, nil
}
