// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartListLogic {
	return &CartListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartListLogic) CartList() (resp []types.CartVO, err error) {
	reply, err := l.svcCtx.TradeRpc.CartList(l.ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	resp = make([]types.CartVO, 0)
	for _, item := range reply.Items {
		resp = append(resp, types.CartVO{
			Id:         item.Id,
			CourseId:   item.CourseId,
			CourseName: item.CourseName,
			CoverUrl:   item.CoverUrl,
			Price:      item.Price,
			NowPrice:   item.NowPrice,
			Expired:    item.Expired,
		})
	}
	return resp, nil
}
