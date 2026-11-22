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
	rpcResp, err := backendlogic.NewSetTodayDataLogic(l.ctx, l.svcCtx.Backend).SetTodayData(&dataclient.TodayDataSetReq{
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
