// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	backendlogic "tjxt/app/media/api/backend/logic"
	mediaclient "tjxt/app/media/api/backend/pb"
	"tjxt/app/media/api/internal/svc"
	"tjxt/app/media/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type FileDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewFileDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FileDeleteLogic {
	return &FileDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *FileDeleteLogic) FileDelete(req *types.FileIdPathReq) (resp *types.OkVO, err error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, err
	}
	if _, err := backendlogic.NewFileDeleteLogic(l.ctx, l.svcCtx.Backend).FileDelete(&mediaclient.FileIdRequest{Id: req.Id}); err != nil {
		return nil, err
	}
	return &types.OkVO{Success: true}, nil
}
