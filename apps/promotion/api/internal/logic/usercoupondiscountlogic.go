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

type UserCouponDiscountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponDiscountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponDiscountLogic {
	return &UserCouponDiscountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponDiscountLogic) UserCouponDiscount(req *types.OrderCouponReq) (resp *types.CouponDiscountVO, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	courses := make([]*pb.OrderCourseDTO, 0, len(req.CourseList))
	for _, c := range req.CourseList {
		courses = append(courses, &pb.OrderCourseDTO{CateId: c.CateId, Id: c.Id, Price: c.Price})
	}
	out, err := l.svcCtx.PromotionRpc.UserCouponDiscount(l.ctx, &pb.OrderCouponDTO{
		UserCouponIds: req.UserCouponIds,
		UserId:        userId,
		CourseList:    courses,
	})
	if err != nil {
		return nil, err
	}
	return toCouponDiscountVO(out), nil
}
