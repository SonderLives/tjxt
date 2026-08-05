package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

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
		return nil, xerr.BadRequestf("业务退款单号不能为空")
	}

	resp, err := l.svcCtx.PayRpc.QueryRefundResult(l.ctx, &payclient.QueryRefundResultRequest{
		BizRefundOrderNo: in.BizRefundOrderId,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询退款结果失败")
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
