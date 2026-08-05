// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package smsplatform

import (
	"context"

	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSmsPlatformsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSmsPlatformsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSmsPlatformsLogic {
	return &ListSmsPlatformsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSmsPlatformsLogic) ListSmsPlatforms() (resp *types.SmsPlatformListVO, err error) {
	// todo: add your logic here and delete this line

	return
}
