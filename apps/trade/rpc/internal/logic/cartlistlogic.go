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

type CartListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartListLogic {
	return &CartListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartListLogic) CartList(in *pb.Empty) (*pb.CartListReply, error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}

	carts, err := l.svcCtx.CartModel.ListByUserId(l.ctx, userId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询购物车失败")
	}

	items := make([]*pb.CartVO, 0, len(carts))
	for _, c := range carts {
		items = append(items, toCartVO(c))
	}
	return &pb.CartListReply{Items: items}, nil
}
