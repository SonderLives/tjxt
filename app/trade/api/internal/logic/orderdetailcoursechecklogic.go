// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/app/trade/api/backend/logic"
	"tjxt/app/trade/api/backend/pb"
	"tjxt/app/trade/api/internal/svc"
	"tjxt/app/trade/api/internal/types"

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
	reply, err := backendlogic.NewOrderDetailCourseCheckLogic(l.ctx, l.svcCtx.Backend).OrderDetailCourseCheck(&pb.IdRequest{Id: req.Id})
	if err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: reply.Value, Message: "ok"}, nil
}
