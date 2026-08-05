// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
	mediaclient "tjxt/apps/media/rpc/media"
	"tjxt/pkg/auth"

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
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.MediaRpc.MediaGet(l.ctx, &mediaclient.MediaIdRequest{MediaId: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.MediaVO{
		Id:         rpcResp.Id,
		Filename:   rpcResp.Filename,
		MediaUrl:   rpcResp.MediaUrl,
		CoverUrl:   rpcResp.CoverUrl,
		Duration:   rpcResp.Duration,
		Size:       rpcResp.Size,
		Status:     rpcResp.Status,
		Creater:    rpcResp.Creater,
		CreateTime: rpcResp.CreateTime,
		UseTimes:   rpcResp.UseTimes,
	}, nil
}
