// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"context"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRoleLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleLogic {
	return &SaveRoleLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveRoleLogic) SaveRole(req *types.RoleSaveReq) (resp *types.IdVO, err error) {
	// todo: add your logic here and delete this line

	return
}
