// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
