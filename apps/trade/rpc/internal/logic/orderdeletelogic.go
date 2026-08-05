package logic

import (
	"context"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDeleteLogic {
	return &OrderDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDeleteLogic) OrderDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if _, err := auth.UserIdFromCtx(l.ctx); err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}

	// 先删明细，再删主单
	details, err := l.svcCtx.OrderDetailModel.ListByOrderId(l.ctx, in.Id)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}
	for _, d := range details {
		if e := l.svcCtx.OrderDetailModel.Delete(l.ctx, d.Id); e != nil {
			return nil, xerr.Wrap(e, xerr.CodeInternal, "删除订单明细失败")
		}
	}

	if err = l.svcCtx.OrderModel.Delete(l.ctx, in.Id); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除订单失败")
	}

	return &pb.Empty{}, nil
}
