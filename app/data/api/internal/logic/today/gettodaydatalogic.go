// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"context"

	backendlogic "tjxt/service/data/api/backend/logic"
	dataclient "tjxt/service/data/api/backend/pb"
	"tjxt/service/data/api/internal/svc"
	"tjxt/service/data/api/internal/types"

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
	rpcResp, err := backendlogic.NewGetTodayDataLogic(l.ctx, l.svcCtx.Backend).GetTodayData(&dataclient.Empty{})
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
