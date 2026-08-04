package service

import (
	"context"
	"database/sql"
	"time"

	"common/idgen"
	"common/xerr"
	"trade/internal/model"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// RefundService 退款业务接口
type RefundService interface {
	// ApplyRefund 学员发起退款申请。
	ApplyRefund(ctx context.Context, userId int64, req *types.RefundApplyReq) error
	// ApproveRefund 审批退款申请（同意/拒绝）。
	ApproveRefund(ctx context.Context, approverId int64, req *types.ApproveReq) error
	// CancelRefund 学员取消退款申请。
	CancelRefund(ctx context.Context, userId int64, req *types.RefundCancelReq) error
	// GetRefundApply 查询退款申请详情（用户/管理端）。
	GetRefundApply(ctx context.Context, id int64) (*types.RefundApplyVO, error)
	// NextRefundApply 获取最早一条待审批的退款申请（管理端）。
	NextRefundApply(ctx context.Context) (*types.RefundApplyVO, error)
	// PageRefundApplies 分页查询退款申请（管理端）。
	PageRefundApplies(ctx context.Context, cond *RefundQuery) (*types.Page, error)
}

// RefundQuery 退款申请查询条件
type RefundQuery struct {
	Id            int64
	OrderId       int64
	OrderDetailId int64
	Mobile        string
	Status        int64
	StartTime     time.Time
	EndTime       time.Time
	PageNo        int64
	PageSize      int64
	IsAsc         bool
}

type refundService struct {
	refundModel    *model.RefundApplyModel
	orderModel     *model.OrderModel
	detailModel    *model.OrderDetailModel
	eventPublisher EventPublisher
	userClient     UserClient
}

// NewRefundService 创建退款业务服务。
func NewRefundService(refundModel *model.RefundApplyModel, orderModel *model.OrderModel, detailModel *model.OrderDetailModel, publisher EventPublisher, userClient UserClient) RefundService {
	return &refundService{
		refundModel:    refundModel,
		orderModel:     orderModel,
		detailModel:    detailModel,
		eventPublisher: publisher,
		userClient:     userClient,
	}
}

// ApplyRefund 学员发起退款申请。
func (s *refundService) ApplyRefund(ctx context.Context, userId int64, req *types.RefundApplyReq) error {
	if req.OrderDetailId == 0 || req.RefundReason == "" {
		return xerr.BadRequestf("订单明细与退款原因不能为空")
	}

	detail, err := s.detailModel.FindById(ctx, req.OrderDetailId)
	if err == sql.ErrNoRows {
		return xerr.NotFound("订单明细不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}
	if detail.UserId != userId {
		return xerr.Forbidden("无权操作该订单")
	}
	refundStatus := formatNullInt(detail.RefundStatus)
	if !canRefund(detail.Status, refundStatus) {
		return xerr.Conflict("当前订单状态不允许申请退款")
	}

	order, err := s.orderModel.FindById(ctx, detail.OrderId)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询订单失败")
	}

	now := time.Now()
	apply := &model.RefundApply{
		Id:            idgen.NextID(),
		OrderDetailId: detail.Id,
		OrderId:       detail.OrderId,
		PayOrderNo:    order.PayOrderNo,
		UserId:        userId,
		RefundAmount:  detail.RealPayAmount,
		Status:        model.RefundStatusPending,
		RefundReason:  req.RefundReason,
		Message:       "待审批",
		QuestionDesc:  sql.NullString{String: req.QuestionDesc, Valid: req.QuestionDesc != ""},
		CreateTime:    now,
		Creater:       userId,
		Updater:       userId,
	}
	if err := s.refundModel.Insert(ctx, apply); err != nil {
		logx.Errorf("apply refund failed, user=%d detail=%d err=%v", userId, req.OrderDetailId, err)
		return xerr.Wrap(err, xerr.CodeInternal, "申请退款失败")
	}

	// 同步更新明细与订单状态
	if err := s.detailModel.MarkRefundStatus(ctx, detail.Id, model.RefundStatusPending); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新订单明细状态失败")
	}
	if err := s.orderModel.MarkRefunding(ctx, detail.OrderId); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新订单状态失败")
	}
	return nil
}

// ApproveRefund 审批退款申请。
func (s *refundService) ApproveRefund(ctx context.Context, approverId int64, req *types.ApproveReq) error {
	if req.Id == 0 {
		return xerr.BadRequestf("退款申请 id 不能为空")
	}
	apply, err := s.refundModel.FindById(ctx, req.Id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("退款申请不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}
	if apply.Status != model.RefundStatusPending {
		return xerr.Conflict("退款申请已处理")
	}

	switch req.ApproveType {
	case 1: // 同意退款
		// 真实环境：调用 pay 服务发起退款，再根据回调结果更新状态。
		// 本实现模拟支付渠道直接退款成功，并发布退款事件通知 learning 撤回课程。
		now := time.Now()
		if err := s.refundModel.Approve(ctx, apply.Id, model.RefundStatusApproved, approverId, req.ApproveOpinion, req.Remark, &now, &now); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "审批退款失败")
		}
		refundOrderNo := idgen.NextID()
		if err := s.refundModel.MarkRefundDone(ctx, apply.Id, model.RefundStatusSuccess, refundOrderNo, "SIMULATED", ""); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "退款处理失败")
		}
		if err := s.detailModel.MarkRefundStatus(ctx, apply.OrderDetailId, model.RefundStatusSuccess); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新订单明细失败")
		}

		if err := s.eventPublisher.PublishRefund(ctx, apply.OrderId, apply.UserId, []int64{mustCourseID(ctx, s.detailModel, apply.OrderDetailId)}, now); err != nil {
			logx.Errorf("publish refund event failed, apply=%d err=%v", apply.Id, err)
		}
		return nil
	case 0: // 拒绝退款
		now := time.Now()
		if err := s.refundModel.Approve(ctx, apply.Id, model.RefundStatusRejected, approverId, req.ApproveOpinion, req.Remark, &now, nil); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "审批退款失败")
		}
		if err := s.detailModel.MarkRefundStatus(ctx, apply.OrderDetailId, model.RefundStatusRejected); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "更新订单明细失败")
		}
		// 恢复订单状态为已支付
		_ = s.orderModel.UpdateStatus(ctx, apply.OrderId, model.OrderStatusPaid, "用户支付成功")
		return nil
	default:
		return xerr.BadRequestf("非法的审批类型")
	}
}

// CancelRefund 学员取消退款申请。
func (s *refundService) CancelRefund(ctx context.Context, userId int64, req *types.RefundCancelReq) error {
	apply, err := s.refundModel.FindById(ctx, req.Id)
	if err == sql.ErrNoRows {
		return xerr.NotFound("退款申请不存在")
	}
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}
	if apply.UserId != userId {
		return xerr.Forbidden("无权操作该退款申请")
	}
	if apply.Status != model.RefundStatusPending {
		return xerr.Conflict("退款申请已处理，无法取消")
	}

	if err := s.refundModel.UpdateStatus(ctx, apply.Id, model.RefundStatusCancelled, "取消退款", nil, nil); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "取消退款失败")
	}
	if err := s.detailModel.MarkRefundStatus(ctx, apply.OrderDetailId, model.RefundStatusCancelled); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "更新订单明细失败")
	}
	// 恢复订单状态为已支付
	_ = s.orderModel.UpdateStatus(ctx, apply.OrderId, model.OrderStatusPaid, "用户支付成功")
	return nil
}

// GetRefundApply 查询退款申请详情。
func (s *refundService) GetRefundApply(ctx context.Context, id int64) (*types.RefundApplyVO, error) {
	apply, err := s.refundModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("退款申请不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}
	return s.buildRefundApplyVO(ctx, apply)
}

// NextRefundApply 获取最早一条待审批的退款申请。
func (s *refundService) NextRefundApply(ctx context.Context) (*types.RefundApplyVO, error) {
	apply, err := s.refundModel.ListPendingOrder(ctx)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}
	return s.buildRefundApplyVO(ctx, apply)
}

// PageRefundApplies 分页查询退款申请。
func (s *refundService) PageRefundApplies(ctx context.Context, q *RefundQuery) (*types.Page, error) {
	offset, limit := normalizePage(q.PageNo, q.PageSize)
	cond := &model.RefundApplyPageCond{
		Id:            q.Id,
		OrderId:       q.OrderId,
		OrderDetailId: q.OrderDetailId,
		Status:        q.Status,
		StartTime:     q.StartTime,
		EndTime:       q.EndTime,
		Offset:        offset,
		Limit:         limit,
		IsAsc:         q.IsAsc,
	}
	if q.Mobile != "" {
		// 手机号过滤需先查询用户服务定位 userId
		userIDs, err := s.resolveUserIdsByMobile(ctx, q.Mobile)
		if err != nil {
			return nil, err
		}
		if len(userIDs) == 0 {
			return &types.Page{List: []types.RefundApplyPageVO{}, Total: 0, Pages: 0}, nil
		}
		// 简化：仅当匹配到单个用户时精确过滤；多用户按首个处理
		cond.UserId = userIDs[0]
	}

	rows, total, err := s.refundModel.ListPage(ctx, cond)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询退款申请失败")
	}

	userIDs := make(map[int64]struct{}, len(rows))
	approverIDs := make(map[int64]struct{}, len(rows))
	for i := range rows {
		userIDs[rows[i].UserId] = struct{}{}
		if rows[i].Approver.Valid {
			approverIDs[rows[i].Approver.Int64] = struct{}{}
		}
	}
	users := s.fetchUsers(ctx, userIDs)
	approvers := s.fetchUsers(ctx, approverIDs)

	list := make([]types.RefundApplyPageVO, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		vo := types.RefundApplyPageVO{
			Id:                r.Id,
			OrderId:           r.OrderId,
			OrderDetailId:     r.OrderDetailId,
			RefundAmount:      r.RefundAmount,
			Status:            r.Status,
			RefundStatusDesc:  refundStatusDesc(r.Status),
			RefundSuccessTime: formatNullTime(r.FinishTime),
			CreateTime:        r.CreateTime.Format(time.RFC3339),
			ApproveTime:       formatNullTime(r.ApproveTime),
		}
		if u, ok := users[r.UserId]; ok {
			vo.ProposerName = u.Name
			vo.ProposerMobile = u.CellPhone
		}
		if a, ok := approvers[formatNullInt(r.Approver)]; ok {
			vo.ApproverName = a.Name
		}
		list = append(list, vo)
	}
	return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
}

// buildRefundApplyVO 组装退款申请详情 VO。
func (s *refundService) buildRefundApplyVO(ctx context.Context, apply *model.RefundApply) (*types.RefundApplyVO, error) {
	vo := &types.RefundApplyVO{
		Id:             apply.Id,
		OrderId:        apply.OrderId,
		OrderDetailId:  apply.OrderDetailId,
		RefundReason:   apply.RefundReason,
		QuestionDesc:   apply.QuestionDesc.String,
		RefundChannel:  apply.RefundChannel.String,
		RefundOrderNo:  formatNullInt(apply.RefundOrderNo),
		PayOrderNo:     formatNullInt(apply.PayOrderNo),
		Status:         apply.Status,
		ApproveOpinion: apply.ApproveOpinion.String,
		ApproveTime:    formatNullTime(apply.ApproveTime),
		FailedReason:   apply.FailedReason.String,
		Message:        apply.Message,
		Remark:         apply.Remark.String,
		CreateTime:     apply.CreateTime.Format(time.RFC3339),
	}

	detail, err := s.detailModel.FindById(ctx, apply.OrderDetailId)
	if err == nil {
		vo.Name = detail.Name
		vo.Price = detail.Price
		vo.RealPayAmount = detail.RealPayAmount
		vo.DiscountAmount = detail.DiscountAmount
		vo.PayChannel = detail.PayChannel
	}

	order, err := s.orderModel.FindById(ctx, apply.OrderId)
	if err == nil {
		vo.PayOrderNo = formatNullInt(order.PayOrderNo)
		vo.OrderTime = order.CreateTime.Format(time.RFC3339)
		vo.PaySuccessTime = formatNullTime(order.PayTime)
	}

	if apply.UserId > 0 {
		users := s.fetchUsers(ctx, map[int64]struct{}{apply.UserId: {}})
		if u, ok := users[apply.UserId]; ok {
			vo.StudentName = u.Name
			vo.Mobile = u.CellPhone
		}
	}
	return vo, nil
}

// fetchUsers 批量获取用户信息。
func (s *refundService) fetchUsers(ctx context.Context, ids map[int64]struct{}) map[int64]*UserInfo {
	if len(ids) == 0 {
		return map[int64]*UserInfo{}
	}
	list := make([]int64, 0, len(ids))
	for id := range ids {
		if id > 0 {
			list = append(list, id)
		}
	}
	users, err := s.userClient.GetByIds(ctx, list)
	if err != nil {
		logx.Errorf("fetch users failed: %v", err)
		return map[int64]*UserInfo{}
	}
	return users
}

// resolveUserIdsByMobile 根据手机号反查用户 id。
func (s *refundService) resolveUserIdsByMobile(ctx context.Context, mobile string) ([]int64, error) {
	// 用户服务未提供按手机号检索的批量接口时，这里返回空；
	// 接入真实用户服务后可补充 /users/search 内部接口。
	_ = mobile
	return nil, nil
}

// mustCourseID 获取明细对应的课程 id（退款事件需要课程信息）。
func mustCourseID(ctx context.Context, detailModel *model.OrderDetailModel, detailID int64) int64 {
	detail, err := detailModel.FindById(ctx, detailID)
	if err != nil {
		return 0
	}
	return detail.CourseId
}
