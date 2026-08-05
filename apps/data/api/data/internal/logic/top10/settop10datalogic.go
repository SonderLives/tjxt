// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTop10DataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTop10DataLogic {
	return &SetTop10DataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetTop10DataLogic) SetTop10Data(req *types.Top10DataSetReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
