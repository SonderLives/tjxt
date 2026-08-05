package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailCourseCheckLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailCourseCheckLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailCourseCheckLogic {
	return &OrderDetailCourseCheckLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailCourseCheckLogic) OrderDetailCourseCheck(in *pb.IdRequest) (*pb.BoolReply, error) {
	d, err := l.svcCtx.OrderDetailModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单明细不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}

	// 已支付/已完成/已报名视为拥有课程权限
	value := d.Status == DetailStatusPaid || d.Status == DetailStatusFinished || d.Status == DetailStatusEnrolled
	return &pb.BoolReply{Value: value}, nil
}
