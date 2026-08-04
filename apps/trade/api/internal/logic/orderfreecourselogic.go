// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderFreeCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderFreeCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderFreeCourseLogic {
	return &OrderFreeCourseLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderFreeCourseLogic) OrderFreeCourse(req *types.FreeCourseReq) (resp *types.PlaceOrderResultVO, err error) {
	// todo: add your logic here and delete this line

	return
}
