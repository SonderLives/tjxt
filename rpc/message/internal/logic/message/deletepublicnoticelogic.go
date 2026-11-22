package messagelogic

import (
	"context"

	"tjxt/rpc/message/internal/svc"
	"tjxt/rpc/message/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除公告失败")
	}
	return &pb.Empty{}, nil
}
