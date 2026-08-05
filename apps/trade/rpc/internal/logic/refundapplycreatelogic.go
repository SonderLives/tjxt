package logic

import (
	"context"
	"database/sql"
	"errors"

	"tjxt/apps/trade/rpc/internal/model"
	"tjxt/apps/trade/rpc/internal/svc"
	"tjxt/apps/trade/rpc/pb"
	"tjxt/pkg/auth"
	"tjxt/pkg/xerr"

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
	userId, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, xerr.New(xerr.CodeUnauthorized, "未登录")
	}
	if in.OrderDetailId <= 0 {
		return nil, xerr.BadRequestf("订单明细ID不能为空")
	}

	detail, err := l.svcCtx.OrderDetailModel.FindOne(l.ctx, in.OrderDetailId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("订单明细不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
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
		return nil, xerr.Wrap(err, xerr.CodeInternal, "创建退款申请失败")
	}

	if err = l.svcCtx.OrderDetailModel.UpdateRefundStatus(l.ctx, in.OrderDetailId, RefundStatusPending); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新订单明细退款状态失败")
	}
	return &pb.Empty{}, nil
}
