// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"context"

	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"

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
	// todo: add your logic here and delete this line

	return
}
