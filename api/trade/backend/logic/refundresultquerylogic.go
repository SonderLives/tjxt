package tradelogic

import (
	"context"

	payclient "tjxt/rpc/pay/client/pay"
	"tjxt/api/trade/backend/svc"
	"tjxt/api/trade/backend/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundResultQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundResultQueryLogic {
	return &RefundResultQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundResultQueryLogic) RefundResultQuery(in *pb.RefundResultQueryRequest) (*pb.RefundResultDTO, error) {
	if in.BizRefundOrderId <= 0 {
		return nil, errorx.BadRequestf("业务退款单号不能为空")
	}

	resp, err := l.svcCtx.PayRpc.QueryRefundResult(l.ctx, &payclient.QueryRefundResultRequest{
		BizRefundOrderNo: in.BizRefundOrderId,
	})
	if err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询退款结果失败")
	}

	return &pb.RefundResultDTO{
		BizPayOrderId:    0,
		BizRefundOrderId: resp.BizRefundOrderNo,
		PayOrderNo:       0,
		RefundOrderNo:    resp.RefundOrderNo,
		Status:           resp.Status,
		PayChannel:       "",
		RefundChannel:    "",
	}, nil
}
