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

type CouponIssueLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCouponIssueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponIssueLogic {
	return &CouponIssueLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CouponIssueLogic) CouponIssue(req *types.CouponIssueReq) error {
	_, err := l.svcCtx.PromotionRpc.CouponIssue(l.ctx, &pb.CouponIssueFormDTO{
		Id:             req.Id,
		IssueBeginTime: req.IssueBeginTime,
		IssueEndTime:   req.IssueEndTime,
		TermBeginTime:  req.TermBeginTime,
		TermDays:       req.TermDays,
		TermEndTime:    req.TermEndTime,
	})
	return err
}
