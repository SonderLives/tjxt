// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"
	"strconv"
	"strings"

	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartBatchDeleteLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCartBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartBatchDeleteLogic {
	return &CartBatchDeleteLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartBatchDeleteLogic) CartBatchDelete(req *types.CartBatchDeleteReq) (resp *types.NamePlaceVO, err error) {
	ids := make([]int64, 0)
	for _, s := range strings.Split(req.Ids, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		id, parseErr := strconv.ParseInt(s, 10, 64)
		if parseErr != nil {
			return nil, xerr.BadRequestf("购物车id格式错误: %s", s)
		}
		ids = append(ids, id)
	}
	if len(ids) == 0 {
		return nil, xerr.BadRequestf("购物车id不能为空")
	}

	if _, err = l.svcCtx.TradeRpc.CartBatchDelete(l.ctx, &pb.CartBatchDeleteRequest{Ids: ids}); err != nil {
		return nil, err
	}
	return &types.NamePlaceVO{Existed: true, Message: "ok"}, nil
}
