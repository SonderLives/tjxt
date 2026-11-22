package tradelogic

import (
	"context"

	"tjxt/api/trade/backend/svc"
	"tjxt/api/trade/backend/pb"
	"tjxt/pkg/errorx"

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
		return nil, errorx.BadRequestf("课程ID不能为空")
	}

	enrollNum, realPayAmount, refundNum, err := l.svcCtx.OrderDetailModel.StatByCourseId(l.ctx, in.CourseId)
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询课程购买信息失败")
	}

	return &pb.PurchaseInfoReply{
		EnrollNum:     enrollNum,
		RealPayAmount: realPayAmount,
		RefundNum:     refundNum,
	}, nil
}
