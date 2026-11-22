package tradelogic

import (
	"context"
	"errors"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/internal/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartDeleteLogic {
	return &CartDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartDeleteLogic) CartDelete(in *pb.IdRequest) (*pb.Empty, error) {
	userId, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("购物车条目ID不能为空")
	}

	c, err := l.svcCtx.CartModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("购物车条目不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询购物车条目失败")
	}
	if c.UserId != userId {
		return nil, errorx.NotFound("购物车条目不存在")
	}

	if err = l.svcCtx.CartModel.Delete(l.ctx, in.Id); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "删除购物车条目失败")
	}
	return &pb.Empty{}, nil
}
