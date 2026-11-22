package tradelogic

import (
	"context"
	"errors"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/internal/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartGetLogic {
	return &CartGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartGetLogic) CartGet(in *pb.IdRequest) (*pb.CartVO, error) {
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
	return toCartVO(c), nil
}
