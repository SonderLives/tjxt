// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	dataclient "tjxt/apps/data/rpc/data/client/data"
	"tjxt/apps/data/rpc/data/pb"

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
	data := make([]*pb.Top10DataSetUnit, 0, len(req.Data))
	for _, u := range req.Data {
		data = append(data, &pb.Top10DataSetUnit{
			Category:    u.Category,
			Name:        u.Name,
			NewStuNum:   int32(u.NewStuNum),
			OrderAmount: u.OrderAmount,
		})
	}
	rpcResp, err := l.svcCtx.DataRpc.SetTop10Data(l.ctx, &dataclient.Top10DataSetReq{
		Version: int32(req.Version),
		Data:    data,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: rpcResp.Success}, nil
}
