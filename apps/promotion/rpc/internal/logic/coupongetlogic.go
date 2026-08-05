package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponGetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponGetLogic {
	return &CouponGetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponGet 根据 id 查询优惠券详情。
func (l *CouponGetLogic) CouponGet(in *pb.IdRequest) (*pb.CouponDetailVO, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("优惠券 id 非法")
	}

	data, err := l.svcCtx.CouponModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("优惠券不存在")
		}
		l.Errorf("find coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询优惠券失败")
	}
	if data.Deleted == 1 {
		return nil, xerr.NotFound("优惠券不存在")
	}
	return toCouponDetailVO(data), nil
}
