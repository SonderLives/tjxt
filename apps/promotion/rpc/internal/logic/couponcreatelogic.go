package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponCreateLogic {
	return &CouponCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponCreate 新增优惠券，创建后处于草稿状态，需再调用 CouponIssue 才会开始发放。
func (l *CouponCreateLogic) CouponCreate(in *pb.CouponFormDTO) (*pb.CouponDetailVO, error) {
	if err := validateCouponForm(in); err != nil {
		return nil, err
	}

	data := buildCoupon(in, 0)
	ret, err := l.svcCtx.CouponModel.Insert(l.ctx, data)
	if err != nil {
		l.Errorf("insert coupon failed, name=%s, err=%v", data.Name, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增优惠券失败")
	}

	id, err := ret.LastInsertId()
	if err != nil {
		l.Errorf("get coupon last insert id failed, err=%v", err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增优惠券失败")
	}

	saved, err := l.svcCtx.CouponModel.FindOne(l.ctx, id)
	if err != nil {
		l.Errorf("find coupon after insert failed, id=%d, err=%v", id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "新增优惠券失败")
	}
	return toCouponDetailVO(saved), nil
}
