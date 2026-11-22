package paylogic

import (
	"context"

	"tjxt/rpc/pay/internal/svc"
	"tjxt/rpc/pay/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type QueryPayResultLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewQueryPayResultLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryPayResultLogic {
	return &QueryPayResultLogic{ctx: ctx, svcCtx: svcCtx, Logger: logx.WithContext(ctx)}
}

func (l *QueryPayResultLogic) QueryPayResult(in *pb.QueryPayResultRequest) (*pb.PayResultResponse, error) {
	if in.BizOrderNo <= 0 {
		return nil, errorx.BadRequestf("biz_order_no 非法")
	}
	m, err := l.svcCtx.PayOrderModel.FindOneByBizOrderNo(l.ctx, in.BizOrderNo)
	if err != nil {
		if isNotFound(err) {
			return nil, errorx.NotFound("支付单不存在")
		}
		return nil, errorx.Wrapf(err, errorx.CodeInternal, "查询支付结果失败")
	}
	return &pb.PayResultResponse{
		PayOrderNo: m.PayOrderNo,
		BizOrderNo: m.BizOrderNo,
		Status:     int32(m.Status),
	}, nil
}
