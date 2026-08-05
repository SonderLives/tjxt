// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetBoardDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetBoardDataLogic {
	return &GetBoardDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetBoardDataLogic) GetBoardData(req *types.BoardDataReq) (resp *types.EchartsVO, err error) {
	// todo: add your logic here and delete this line

	return
}
