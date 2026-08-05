// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMediaListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaListLogic {
	return &MediaListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MediaListLogic) MediaList(req *types.MediaListReq) (resp *types.MediaListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
