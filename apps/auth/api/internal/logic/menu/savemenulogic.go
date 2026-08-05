// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMenuLogic {
	return &SaveMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveMenuLogic) SaveMenu(req *types.MenuSaveReq) (resp *types.IdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
