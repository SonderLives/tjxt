// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"context"

	"tjxt/api/data/internal/svc"
	"tjxt/api/data/internal/types"
	dataclient "tjxt/rpc/data/client/data"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTodayDataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTodayDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTodayDataLogic {
	return &GetTodayDataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTodayDataLogic) GetTodayData() (resp *types.TodayDataVO, err error) {
	rpcResp, err := l.svcCtx.DataRpc.GetTodayData(l.ctx, &dataclient.Empty{})
	if err != nil {
		return nil, err
	}
	return &types.TodayDataVO{
		Visits:      rpcResp.Visits,
		OrderAmount: rpcResp.OrderAmount,
		OrderNum:    int(rpcResp.OrderNum),
		StuNewNum:   int(rpcResp.StuNewNum),
	}, nil
}
