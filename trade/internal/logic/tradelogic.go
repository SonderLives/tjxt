// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type TradeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewTradeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *TradeLogic {
	return &TradeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *TradeLogic) Trade(req *types.Request) (resp *types.Response, err error) {
	// todo: add your logic here and delete this line

	return
}
