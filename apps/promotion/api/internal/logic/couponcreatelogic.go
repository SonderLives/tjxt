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

type CouponCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponCreateLogic {
	return &CouponCreateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponCreateLogic) CouponCreate(req *types.CouponFormReq) (resp *types.CouponDetailVO, err error) {
	out, err := l.svcCtx.PromotionRpc.CouponCreate(l.ctx, &pb.CouponFormDTO{
		Id:                req.Id,
		Name:              req.Name,
		DiscountType:      req.DiscountType,
		DiscountValue:     req.DiscountValue,
		MaxDiscountAmount: req.MaxDiscountAmount,
		ThresholdAmount:   req.ThresholdAmount,
		ObtainWay:         req.ObtainWay,
		Specific:          req.Specific,
		Scopes:            req.Scopes,
		TotalNum:          req.TotalNum,
		UserLimit:         req.UserLimit,
	})
	if err != nil {
		return nil, err
	}
	return toCouponDetailVO(out), nil
}
