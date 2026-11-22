// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"context"

	backendlogic "tjxt/api/data/backend/logic"
	dataclient "tjxt/api/data/backend/pb"
	"tjxt/api/data/internal/svc"
	"tjxt/api/data/internal/types"

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
	rpcResp, err := backendlogic.NewSetBoardDataLogic(l.ctx, l.svcCtx.Backend).SetBoardData(&dataclient.BoardDataSetReq{
		Version: int32(req.Version),
		Type:    int32(req.Type),
		Data:    req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: rpcResp.Success}, nil
}
