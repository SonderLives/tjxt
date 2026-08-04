// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailCourseCheckLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailCourseCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailCourseCheckLogic {
	return &OrderDetailCourseCheckLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *OrderDetailCourseCheckLogic) OrderDetailCourseCheck(req *types.OrderIdReq) (resp *types.NamePlaceVO, err error) {
	// todo: add your logic here and delete this line

	return
}
