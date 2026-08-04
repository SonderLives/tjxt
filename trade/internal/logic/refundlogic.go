package logic

import (
	"context"

	"common/auth"
	"common/result"
	"trade/internal/service"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 发起退款申请 ============

type RefundApplyCreateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCreateLogic {
	return &RefundApplyCreateLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyCreateLogic) RefundApplyCreate(req *types.RefundApplyReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.RefundService.ApplyRefund(l.ctx, userID, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 审批退款 ============

type RefundApplyApproveLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyApproveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyApproveLogic {
	return &RefundApplyApproveLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyApproveLogic) RefundApplyApprove(req *types.ApproveReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.RefundService.ApproveRefund(l.ctx, userID, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 取消退款申请 ============

type RefundApplyCancelLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyCancelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyCancelLogic {
	return &RefundApplyCancelLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyCancelLogic) RefundApplyCancel(req *types.RefundCancelReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	if err := l.svcCtx.RefundService.CancelRefund(l.ctx, userID, req); err != nil {
		return nil, err
	}
	return result.Ok(), nil
}

// ============ 退款申请详情 ============

type RefundApplyDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyDetailLogic {
	return &RefundApplyDetailLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyDetailLogic) RefundApplyDetail(req *types.RefundIdReq) (resp *result.R, err error) {
	apply, err := l.svcCtx.RefundService.GetRefundApply(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(apply), nil
}

// ============ 下一条待审批退款 ============

type RefundApplyNextLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyNextLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyNextLogic {
	return &RefundApplyNextLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyNextLogic) RefundApplyNext() (resp *result.R, err error) {
	apply, err := l.svcCtx.RefundService.NextRefundApply(l.ctx)
	if err != nil {
		return nil, err
	}
	return result.OkData(apply), nil
}

// ============ 退款申请分页 ============

type RefundApplyPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyPageLogic {
	return &RefundApplyPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyPageLogic) RefundApplyPage(req *types.RefundApplyPageReq) (resp *result.R, err error) {
	q := &service.RefundQuery{
		Id:            req.Id,
		OrderId:       req.OrderId,
		OrderDetailId: req.OrderDetailId,
		Mobile:        req.Mobile,
		Status:        req.RefundStatus,
		PageNo:        req.PageNo,
		PageSize:      req.PageSize,
		IsAsc:         req.IsAsc,
	}
	if t, ok := parseTime(req.ApplyStartTime); ok {
		q.StartTime = t
	}
	if t, ok := parseTime(req.ApplyEndTime); ok {
		q.EndTime = t
	}
	page, err := l.svcCtx.RefundService.PageRefundApplies(l.ctx, q)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 退款申请（按 id 查询，兼容） ============

type RefundApplyGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundApplyGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundApplyGetLogic {
	return &RefundApplyGetLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *RefundApplyGetLogic) RefundApplyGet(req *types.RefundIdReq) (resp *result.R, err error) {
	apply, err := l.svcCtx.RefundService.GetRefundApply(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(apply), nil
}
