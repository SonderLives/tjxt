package admin

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserByIdLogic {
	return &UpdateUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateUserById 管理端更新指定用户。
func (l *UpdateUserByIdLogic) UpdateUserById(req *types.UserUpdateReq) error {
	_, err := l.svcCtx.UserRpc.UpdateUserById(l.ctx, convert.FromUserUpdateReq(req))
	return err
}
