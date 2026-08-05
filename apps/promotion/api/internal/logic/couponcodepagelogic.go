// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/promotion/api/internal/svc"
	"tjxt/apps/promotion/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/apps/promotion/rpc/promotion"
)

type CouponCodePageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponCodePageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponCodePageLogic {
	return &CouponCodePageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponCodePageLogic) CouponCodePage(req *types.CouponCodePageReq) (resp *types.CouponCodePageResp, err error) {
	out, err := l.svcCtx.PromotionRpc.CouponCodePage(l.ctx, &pb.CouponCodePageRequest{
		Page:     toPageRequest(req.PageRequest),
		CouponId: req.CouponId,
		Status:   req.Status,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.ExchangeCodeVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toExchangeCodeVO(v))
	}
	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}
	pages := int64(0)
	if limit > 0 {
		pages = (out.Total + limit - 1) / limit
	}
	return &types.CouponCodePageResp{List: list, Total: out.Total, Pages: pages}, nil
}
