package logic

import (
	"context"

	"tjxt/apps/message/rpc/internal/svc"
	"tjxt/apps/message/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteNoticeTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteNoticeTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteNoticeTaskLogic {
	return &DeleteNoticeTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 通知任务 物理删除
func (l *DeleteNoticeTaskLogic) DeleteNoticeTask(in *pb.IdReq) (*pb.Empty, error) {
	if err := l.svcCtx.NoticeTaskModel.Delete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除通知任务失败")
	}
	return &pb.Empty{}, nil
}
