// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	backendlogic "tjxt/service/media/api/backend/logic"
	mediaclient "tjxt/service/media/api/backend/pb"
	"tjxt/service/media/api/internal/svc"
	"tjxt/service/media/api/internal/types"
	"tjxt/common/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileSaveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileSaveLogic {
	return &FileSaveLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileSaveLogic) FileSave(req *types.FileSaveReq) (resp *types.FileIdVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	rpcResp, err := backendlogic.NewFileSaveLogic(l.ctx, l.svcCtx.Backend).FileSave(&mediaclient.FileSaveRequest{
		Filename: req.Filename,
		Key:      req.Key,
		Size:     req.Size,
	})
	if err != nil {
		return nil, err
	}
	return &types.FileIdVO{Id: rpcResp.Id}, nil
}
