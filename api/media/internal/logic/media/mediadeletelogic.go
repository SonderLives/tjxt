// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	backendlogic "tjxt/api/media/backend/logic"
	mediaclient "tjxt/api/media/backend/pb"
	"tjxt/api/media/internal/svc"
	"tjxt/api/media/internal/types"
	"tjxt/pkg/auth"

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
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if _, err := backendlogic.NewMediaDeleteLogic(l.ctx, l.svcCtx.Backend).MediaDelete(&mediaclient.MediaIdRequest{MediaId: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
