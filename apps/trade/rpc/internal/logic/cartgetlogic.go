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
		return nil, xerr.BadRequestf("购物车条目ID不能为空")
	}

	c, err := l.svcCtx.CartModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("购物车条目不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询购物车条目失败")
	}
	return toCartVO(c), nil
}
