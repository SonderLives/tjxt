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

type CouponPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPageLogic {
	return &CouponPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponPageLogic) CouponPage(req *types.CouponPageReq) (resp *types.CouponPageResp, err error) {
	out, err := l.svcCtx.PromotionRpc.CouponPage(l.ctx, &pb.CouponPageRequest{
		Page:   toPageRequest(req.PageRequest),
		Name:   req.Name,
		Status: req.Status,
		Type:   req.Type,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CouponPageVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toCouponPageVO(v))
	}
	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}
	pages := int64(0)
	if limit > 0 {
		pages = (out.Total + limit - 1) / limit
	}
	return &types.CouponPageResp{List: list, Total: out.Total, Pages: pages}, nil
}
