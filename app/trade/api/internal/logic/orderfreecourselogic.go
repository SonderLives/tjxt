// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	backendlogic "tjxt/service/trade/api/backend/logic"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"

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
	reply, err := backendlogic.NewOrderFreeCourseLogic(l.ctx, l.svcCtx.Backend).OrderFreeCourse(&pb.FreeCourseRequest{CourseId: req.CourseId})
	if err != nil {
		return nil, err
	}

	return &types.PlaceOrderResultVO{
		OrderId:    reply.OrderId,
		PayAmount:  reply.PayAmount,
		PayOutTime: reply.PayOutTime,
		Status:     int64(reply.Status),
	}, nil
}
