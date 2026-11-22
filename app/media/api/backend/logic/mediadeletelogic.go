package medialogic

import (
	"context"

	"tjxt/service/media/api/backend/pb"
	"tjxt/service/media/api/backend/svc"
	"tjxt/common/errorx"

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
		return nil, errorx.BadRequestf("mediaId 不能为空")
	}
	// 逻辑删除（deleted=1），删除不存在的 id 幂等返回成功
	if err := l.svcCtx.MediaModel.SoftDelete(l.ctx, in.MediaId); err != nil {
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "删除媒资失败")
	}
	return &pb.Empty{}, nil
}
