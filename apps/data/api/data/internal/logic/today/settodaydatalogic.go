// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	dataclient "tjxt/apps/data/rpc/data/client/data"

	"github.com/zeromicro/go-zero/core/logx"
)

type SetTodayDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetTodayDataLogic {
	return &SetTodayDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetTodayDataLogic) SetTodayData(req *types.TodayDataSetReq) (resp *types.OkVO, err error) {
	rpcResp, err := l.svcCtx.DataRpc.SetTodayData(l.ctx, &dataclient.TodayDataSetReq{
		Version:     int32(req.Version),
		Visits:      req.Visits,
		OrderAmount: req.OrderAmount,
		OrderNum:    int32(req.OrderNum),
		StuNewNum:   int32(req.StuNewNum),
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: rpcResp.Success}, nil
}
