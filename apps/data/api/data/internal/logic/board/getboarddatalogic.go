// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	dataclient "tjxt/apps/data/rpc/data/client/data"

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
	typesArr := make([]int32, 0, len(req.Types))
	for _, t := range req.Types {
		typesArr = append(typesArr, int32(t))
	}
	rpcResp, err := l.svcCtx.DataRpc.GetBoardData(l.ctx, &dataclient.BoardDataReq{Types: typesArr})
	if err != nil {
		return nil, err
	}

	resp = &types.EchartsVO{}
	for _, x := range rpcResp.XAxis {
		resp.XAxis = append(resp.XAxis, types.AxisVO{
			Type:     x.Type,
			Max:      x.Max,
			Min:      x.Min,
			Average:  x.Average,
			Data:     x.Data,
			Interval: x.Interval,
		})
	}
	for _, y := range rpcResp.YAxis {
		resp.YAxis = append(resp.YAxis, types.AxisVO{
			Type:     y.Type,
			Max:      y.Max,
			Min:      y.Min,
			Average:  y.Average,
			Data:     y.Data,
			Interval: y.Interval,
		})
	}
	for _, s := range rpcResp.Series {
		resp.Series = append(resp.Series, types.SerierVO{
			Name: s.Name,
			Type: s.Type,
			Data: s.Data,
			Max:  s.Max,
			Min:  s.Min,
		})
	}
	return
}
