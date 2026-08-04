// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayChannelGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPayChannelGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayChannelGetLogic {
	return &PayChannelGetLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PayChannelGetLogic) PayChannelGet(req *types.PayChannelIdReq) (resp *types.PayChannelDTO, err error) {
	// todo: add your logic here and delete this line

	return
}
