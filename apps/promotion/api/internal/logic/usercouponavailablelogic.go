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

type UserCouponAvailableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUserCouponAvailableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserCouponAvailableLogic {
	return &UserCouponAvailableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserCouponAvailableLogic) UserCouponAvailable(req *types.UserCouponAvailableReq) (resp *types.CouponDiscountListResp, err error) {
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	courses := make([]*pb.OrderCourseDTO, 0, len(req.CourseList))
	for _, c := range req.CourseList {
		courses = append(courses, &pb.OrderCourseDTO{CateId: c.CateId, Id: c.Id, Price: c.Price})
	}
	out, err := l.svcCtx.PromotionRpc.UserCouponAvailable(l.ctx, &pb.OrderCourseListRequest{
		UserId:     userId,
		CourseList: courses,
	})
	if err != nil {
		return nil, err
	}
	list := make([]types.CouponDiscountVO, 0, len(out.List))
	for _, v := range out.List {
		list = append(list, *toCouponDiscountVO(v))
	}
	return &types.CouponDiscountListResp{List: list}, nil
}
