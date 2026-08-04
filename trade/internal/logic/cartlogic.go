package logic

import (
	"context"

	"common/auth"
	"common/result"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 购物车列表 ============

type CartListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartListLogic {
	return &CartListLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CartListLogic) CartList() (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	list, err := l.svcCtx.CartService.ListCart(l.ctx, userID)
	if err != nil {
		return nil, err
	}
	return result.OkData(list), nil
}

// ============ 添加购物车 ============

type CartAddLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartAddLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartAddLogic {
	return &CartAddLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CartAddLogic) CartAdd(req *types.CartsAddReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.CartService.AddToCart(l.ctx, userID, req.CourseId); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 删除购物车 ============

type CartDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartDeleteLogic {
	return &CartDeleteLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *CartDeleteLogic) CartDelete(req *types.CartsDeleteReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	ids := parseIDs(req.Ids)
	if len(ids) == 0 {
		return nil, errBadRequest("无效的购物车条目")
	}
	if err := l.svcCtx.CartService.DeleteCart(l.ctx, userID, ids); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}
