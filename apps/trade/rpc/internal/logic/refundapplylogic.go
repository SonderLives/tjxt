package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyLogic {
	return &RefundApplyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 退款 =====
func (l *RefundApplyLogic) RefundApply(in *pb.RefundApplyRequest) (*pb.RefundResultDTO, error) {
	if in.BizOrderNo <= 0 {
		return nil, xerr.BadRequestf("业务订单号不能为空")
	}

	resp, err := l.svcCtx.PayRpc.ApplyRefund(l.ctx, &payclient.ApplyRefundRequest{
		BizOrderNo:       in.BizOrderNo,
		BizRefundOrderNo: in.BizRefundOrderNo,
		RefundAmount:     in.RefundAmount,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "申请退款失败")
	}

	return &pb.RefundResultDTO{
		BizPayOrderId:    in.BizOrderNo,
		BizRefundOrderId: resp.BizRefundOrderNo,
		PayOrderNo:       0,
		RefundOrderNo:    resp.RefundOrderNo,
		Status:           resp.Status,
		PayChannel:       "",
		RefundChannel:    "",
	}, nil
}
