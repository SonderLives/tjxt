// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetBoardDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetBoardDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetBoardDataLogic {
	return &SetBoardDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetBoardDataLogic) SetBoardData(req *types.BoardDataSetReq) (resp *types.OkVO, err error) {
	// todo: add your logic here and delete this line

	return
}
