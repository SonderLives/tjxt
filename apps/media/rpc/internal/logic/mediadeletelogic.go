package logic

import (
	"context"

	"tjxt/apps/media/rpc/internal/svc"
	"tjxt/apps/media/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaDeleteLogic {
	return &MediaDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *MediaDeleteLogic) MediaDelete(in *pb.MediaIdRequest) (*pb.Empty, error) {
	if in.MediaId <= 0 {
		return nil, xerr.BadRequestf("mediaId 不能为空")
	}
	// 逻辑删除（deleted=1），删除不存在的 id 幂等返回成功
	if err := l.svcCtx.MediaModel.SoftDelete(l.ctx, in.MediaId); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeInternal, "删除媒资失败")
	}
	return &pb.Empty{}, nil
}
