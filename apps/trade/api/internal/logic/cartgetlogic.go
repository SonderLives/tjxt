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

type CartGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartGetLogic {
	return &CartGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartGetLogic) CartGet(req *types.CartIdReq) (resp *types.CartVO, err error) {
	reply, err := l.svcCtx.TradeRpc.CartGet(l.ctx, &pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}

	return &types.CartVO{
		Id:         reply.Id,
		CourseId:   reply.CourseId,
		CourseName: reply.CourseName,
		CoverUrl:   reply.CoverUrl,
		Price:      reply.Price,
		NowPrice:   reply.NowPrice,
		Expired:    reply.Expired,
	}, nil
}
