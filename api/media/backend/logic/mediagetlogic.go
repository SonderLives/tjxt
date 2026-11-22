package medialogic

import (
	"context"
	"errors"

	"tjxt/api/media/backend/model"
	"tjxt/api/media/backend/pb"
	"tjxt/api/media/backend/svc"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewMediaGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaGetLogic {
	return &MediaGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 媒资管理
func (l *MediaGetLogic) MediaGet(in *pb.MediaIdRequest) (*pb.MediaVO, error) {
	if in.MediaId <= 0 {
		return nil, errorx.BadRequestf("mediaId 不能为空")
	}
	media, err := l.svcCtx.MediaModel.FindOneNotDeleted(l.ctx, in.MediaId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("媒资不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询媒资失败")
	}
	return toMediaVO(media), nil
}
