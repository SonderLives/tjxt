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

type MediaSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMediaSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *MediaSaveLogic {
	return &MediaSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MediaSaveLogic) MediaSave(req *types.MediaSaveReq) (resp *types.MediaIdVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	// RPC 协议暂无用户字段，此处仅校验登录态
	rpcResp, err := backendlogic.NewMediaSaveLogic(l.ctx, l.svcCtx.Backend).MediaSave(&mediaclient.MediaSaveRequest{
		Id:       req.Id,
		Filename: req.Filename,
		Duration: req.Duration,
		Size:     req.Size,
		FileId:   req.FileId,
	})
	if err != nil {
		return nil, err
	}
	return &types.MediaIdVO{Id: rpcResp.Id}, nil
}
