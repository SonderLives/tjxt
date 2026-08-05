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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
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
	if c.UserId != userId {
		return nil, xerr.NotFound("购物车条目不存在")
	}

	if err = l.svcCtx.CartModel.Delete(l.ctx, in.Id); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除购物车条目失败")
	}
	return &pb.Empty{}, nil
}
