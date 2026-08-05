package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponPauseLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponPauseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponPauseLogic {
	return &CouponPauseLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponPause 暂停发放优惠券，仅发放中的券可暂停；已领取的券不受影响。
func (l *CouponPauseLogic) CouponPause(in *pb.IdRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("优惠券 id 非法")
	}

	data, err := l.svcCtx.CouponModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("优惠券不存在")
		}
		l.Errorf("find coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "暂停发放失败")
	}
	if data.Deleted == 1 {
		return nil, xerr.NotFound("优惠券不存在")
	}
	if data.Status != model.CouponStatusIssued {
		return nil, xerr.Conflict("只有发放中的优惠券才能暂停")
	}

	if err := l.svcCtx.CouponModel.UpdateStatus(l.ctx, in.Id, model.CouponStatusPaused, in.UserId); err != nil {
		l.Errorf("pause coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "暂停发放失败")
	}
	return &pb.Empty{}, nil
}
