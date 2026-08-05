// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTop10DataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTop10DataLogic {
	return &GetTop10DataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTop10DataLogic) GetTop10Data() (resp *types.Top10DataVO, err error) {
	// todo: add your logic here and delete this line

	return
}
