package tradelogic

import (
	"context"
	"errors"

	"tjxt/service/trade/api/backend/model"
	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/auth"
	"tjxt/common/errorx"

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
		return nil, err
	}

	carts, err := l.svcCtx.CartModel.ListByUserId(l.ctx, userId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询购物车失败")
	}

	items := make([]*pb.CartVO, 0, len(carts))
	for _, c := range carts {
		items = append(items, toCartVO(c))
	}
	return &pb.CartListReply{Items: items}, nil
}
