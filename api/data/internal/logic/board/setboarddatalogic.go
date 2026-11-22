// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"context"

	"tjxt/api/data/internal/svc"
	"tjxt/api/data/internal/types"
	dataclient "tjxt/rpc/data/client/data"

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
	rpcResp, err := l.svcCtx.DataRpc.SetBoardData(l.ctx, &dataclient.BoardDataSetReq{
		Version: int32(req.Version),
		Type:    int32(req.Type),
		Data:    req.Data,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: rpcResp.Success}, nil
}
