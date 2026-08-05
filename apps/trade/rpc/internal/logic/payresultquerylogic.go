package logic

import (
	"context"

	payclient "tjxt/apps/pay/rpc/pay"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayResultQueryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayResultQueryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayResultQueryLogic {
	return &PayResultQueryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayResultQueryLogic) PayResultQuery(in *pb.PayResultQueryRequest) (*pb.PayResultDTO, error) {
	if in.BizOrderId <= 0 {
		return nil, xerr.BadRequestf("业务订单ID不能为空")
	}

	resp, err := l.svcCtx.PayRpc.QueryPayResult(l.ctx, &payclient.QueryPayResultRequest{
		BizOrderNo: in.BizOrderId,
	})
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询支付结果失败")
	}

	// pay 侧未返回支付渠道与支付成功时间，留空
	return &pb.PayResultDTO{
		BizOrderId: in.BizOrderId,
		Status:     int32(resp.Status),
		PayOrderNo: resp.PayOrderNo,
	}, nil
}
