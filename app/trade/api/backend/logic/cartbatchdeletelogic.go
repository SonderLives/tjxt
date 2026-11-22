package tradelogic

import (
	"context"
	"errors"

	"tjxt/service/trade/api/backend/model"
	"tjxt/service/trade/api/backend/svc"
	"tjxt/service/trade/api/backend/pb"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type CartBatchDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCartBatchDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CartBatchDeleteLogic {
	return &CartBatchDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CartBatchDeleteLogic) CartBatchDelete(in *pb.CartBatchDeleteRequest) (*pb.Empty, error) {
	for _, id := range in.Ids {
		if id <= 0 {
			continue
		}
		// 条目不存在视为已删除，跳过
		if err := l.svcCtx.CartModel.Delete(l.ctx, id); err != nil && !errors.Is(err, model.ErrNotFound) {
			return nil, errorx.Wrap(err, errorx.CodeInternal, "批量删除购物车条目失败")
		}
	}
	return &pb.Empty{}, nil
}
