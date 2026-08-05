package logic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyApproveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyApproveLogic {
	return &RefundApplyApproveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundApplyApproveLogic) RefundApplyApprove(in *pb.ApproveRequest) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("退款申请ID不能为空")
	}

	ra, err := l.svcCtx.RefundApplyModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("退款申请不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}

	// approve_type：1 同意 2 拒绝
	status := RefundStatusReject
	if in.ApproveType == 1 {
		status = RefundStatusApprove
	}

	if err = l.svcCtx.RefundApplyModel.UpdateApprove(l.ctx, in.Id, status, 0, in.ApproveOpinion, in.Remark,
		sql.NullTime{Time: now(), Valid: true}, 0, 0); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "审批退款申请失败")
	}

	if err = l.svcCtx.OrderDetailModel.UpdateRefundStatus(l.ctx, ra.OrderDetailId, status); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新订单明细退款状态失败")
	}
	return &pb.Empty{}, nil
}
