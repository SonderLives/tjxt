package tradelogic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/rpc/trade/internal/model"
	"tjxt/rpc/trade/internal/svc"
	"tjxt/rpc/trade/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundApplyCreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCreateLogic {
	return &RefundApplyCreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ===== 退款申请 =====
func (l *RefundApplyCreateLogic) RefundApplyCreate(in *pb.RefundApplyFormRequest) (*pb.Empty, error) {
	userId, err := auth.UserIdFromIncomingMD(l.ctx)
	if err != nil {
		return nil, err
	}
	if in.OrderDetailId <= 0 {
		return nil, errorx.BadRequestf("订单明细ID不能为空")
	}

	detail, err := l.svcCtx.OrderDetailModel.FindOne(l.ctx, in.OrderDetailId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("订单明细不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询订单明细失败")
	}

	ra := &model.RefundApply{
		Id:            nextID(),
		UserId:        userId,
		OrderId:       detail.OrderId,
		OrderDetailId: in.OrderDetailId,
		RefundAmount:  detail.RealPayAmount,
		Status:        RefundStatusPending,
		RefundReason:  in.RefundReason,
		Message:       refundStatusDesc(RefundStatusPending),
		QuestionDesc:  sql.NullString{String: in.QuestionDesc, Valid: in.QuestionDesc != ""},
		Creater:       userId,
		Updater:       userId,
		CreateTime:    now(),
	}
	if _, err = l.svcCtx.RefundApplyModel.Insert(l.ctx, ra); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "创建退款申请失败")
	}

	if err = l.svcCtx.OrderDetailModel.UpdateRefundStatus(l.ctx, in.OrderDetailId, RefundStatusPending); err != nil {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "更新订单明细退款状态失败")
	}
	return &pb.Empty{}, nil
}
