package logic

import (
	"context"

	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type OrderDetailPurchaseInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewOrderDetailPurchaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPurchaseInfoLogic {
	return &OrderDetailPurchaseInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *OrderDetailPurchaseInfoLogic) OrderDetailPurchaseInfo(in *pb.PurchaseInfoRequest) (*pb.PurchaseInfoReply, error) {
	if in.CourseId <= 0 {
		return nil, xerr.BadRequestf("课程ID不能为空")
	}

	enrollNum, realPayAmount, refundNum, err := l.svcCtx.OrderDetailModel.StatByCourseId(l.ctx, in.CourseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程购买信息失败")
	}

	return &pb.PurchaseInfoReply{
		EnrollNum:     enrollNum,
		RealPayAmount: realPayAmount,
		RefundNum:     refundNum,
	}, nil
}
