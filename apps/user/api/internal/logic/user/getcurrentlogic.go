// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"context"

	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCurrentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentLogic {
	return &GetCurrentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCurrentLogic) GetCurrent() (resp *types.UserVO, err error) {
	// todo: add your logic here and delete this line

	return
}
