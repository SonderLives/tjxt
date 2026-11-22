package user

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCurrentUserLogic {
	return &UpdateCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateCurrentUser 更新当前登录用户，id 取自 JWT 载荷。
func (l *UpdateCurrentUserLogic) UpdateCurrentUser(req *types.UserFormReq) error {
	uid, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return err
	}
	in := convert.FromUserFormReq(req)
	in.Id = uid
	_, err = l.svcCtx.UserRpc.UpdateCurrentUser(l.ctx, in)
	return err
}
