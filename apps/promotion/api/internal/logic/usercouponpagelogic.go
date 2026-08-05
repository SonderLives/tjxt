// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package logic

import (
	"context"

	"tjxt/apps/promotion/api/internal/svc"
	"tjxt/apps/promotion/api/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/apps/promotion/rpc/promotion"
)

type UserCouponPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponPageLogic {
	return &UserCouponPageLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponPageLogic) UserCouponPage(req *types.UserCouponPageReq) (resp *types.UserCouponPageResp, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.PromotionRpc.UserCouponPage(l.ctx, &pb.UserCouponPageRequest{
		Page:   toPageRequest(req.PageRequest),
		Status: req.Status,
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CouponVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toCouponVO(v))
	}
	limit := req.PageSize
	if limit <= 0 {
		limit = 10
	}
	pages := int64(0)
	if limit > 0 {
		pages = (out.Total + limit - 1) / limit
	}
	return &types.UserCouponPageResp{List: list, Total: out.Total, Pages: pages}, nil
}
