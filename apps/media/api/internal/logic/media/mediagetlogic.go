// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMediaGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaGetLogic {
	return &MediaGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MediaGetLogic) MediaGet(req *types.MediaIdPathReq) (resp *types.MediaVO, err error) {
	// todo: add your logic here and delete this line

	return
}
