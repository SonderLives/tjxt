// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"context"

	backendlogic "tjxt/app/data/api/backend/logic"
	dataclient "tjxt/app/data/api/backend/pb"
	"tjxt/app/data/api/internal/svc"
	"tjxt/app/data/api/internal/types"

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
	data := make([]*dataclient.Top10DataSetUnit, 0, len(req.Data))
	for _, u := range req.Data {
		data = append(data, &dataclient.Top10DataSetUnit{
			Category:    u.Category,
			Name:        u.Name,
			NewStuNum:   int32(u.NewStuNum),
			OrderAmount: u.OrderAmount,
		})
	}
	rpcResp, err := backendlogic.NewSetTop10DataLogic(l.ctx, l.svcCtx.Backend).SetTop10Data(&dataclient.Top10DataSetReq{
		Version: int32(req.Version),
		Data:    data,
	})
	if err != nil {
		return nil, err
	}
	return &types.OkVO{Success: rpcResp.Success}, nil
}
