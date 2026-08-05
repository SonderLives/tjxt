// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"context"

	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	dataclient "tjxt/apps/data/rpc/data/client/data"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTop10DataLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetTop10DataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTop10DataLogic {
	return &GetTop10DataLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTop10DataLogic) GetTop10Data() (resp *types.Top10DataVO, err error) {
	rpcResp, err := l.svcCtx.DataRpc.GetTop10Data(l.ctx, &dataclient.Empty{})
	if err != nil {
		return nil, err
	}

	resp = &types.Top10DataVO{}
	for _, c := range rpcResp.Hot {
		resp.Hot = append(resp.Hot, types.CourseInfo{
			Category:    c.Category,
			Name:        c.Name,
			NewStuNum:   int(c.NewStuNum),
			OrderAmount: c.OrderAmount,
		})
	}
	for _, c := range rpcResp.HotSales {
		resp.HotSales = append(resp.HotSales, types.CourseInfo{
			Category:    c.Category,
			Name:        c.Name,
			NewStuNum:   int(c.NewStuNum),
			OrderAmount: c.OrderAmount,
		})
	}
	return
}
