package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponDeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponDeleteLogic {
	return &CouponDeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponDelete 逻辑删除优惠券，仅允许删除草稿或已暂停的券，避免影响已发放数据。
func (l *CouponDeleteLogic) CouponDelete(in *pb.IdRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("优惠券 id 非法")
	}

	rows, err := l.svcCtx.CouponModel.SoftDelete(l.ctx, in.Id, in.UserId)
	if err != nil {
		l.Errorf("delete coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "删除优惠券失败")
	}
	if rows == 0 {
		return nil, xerr.Conflict("优惠券不存在或当前状态不允许删除")
	}
	return &pb.Empty{}, nil
}
