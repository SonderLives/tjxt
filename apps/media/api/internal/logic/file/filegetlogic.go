// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
	mediaclient "tjxt/apps/media/rpc/media"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileGetLogic {
	return &FileGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileGetLogic) FileGet(req *types.FileIdPathReq) (resp *types.FileVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := l.svcCtx.MediaRpc.FileGet(l.ctx, &mediaclient.FileIdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.FileVO{
		Id:       rpcResp.Id,
		Key:      rpcResp.Key,
		Filename: rpcResp.Filename,
		Path:     rpcResp.Path,
		Status:   rpcResp.Status,
	}, nil
}
