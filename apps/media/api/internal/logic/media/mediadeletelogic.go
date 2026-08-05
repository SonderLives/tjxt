// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type MediaDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMediaDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaDeleteLogic {
	return &MediaDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MediaDeleteLogic) MediaDelete(req *types.MediaIdPathReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
