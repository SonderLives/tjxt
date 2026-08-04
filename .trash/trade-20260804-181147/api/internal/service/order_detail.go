package service

import (
	"context"
	"database/sql"
	"time"

	"tjxt/pkg/xerr"
	"tjxt/apps/trade/api/internal/model"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// OrderDetailQuery 订单明细查询条件（管理端分页）
type OrderDetailQuery struct {
	Id           int64
	Mobile       string
	Status       int64
	RefundStatus int64
	PayChannel   string
	StartTime    time.Time
	EndTime      time.Time
	PageNo       int64
	PageSize     int64
	IsAsc        bool
}

// OrderDetailService 订单明细业务接口
type OrderDetailService interface {
	// IsCoursePurchased 用户是否已购买课程。
	IsCoursePurchased(ctx context.Context, userId, courseId int64) (bool, error)
	// EnrollCourseAmounts 批量查询学员报名课程数。
	EnrollCourseAmounts(ctx context.Context, studentIds []int64) (map[int64]int64, error)
	// EnrollNumByCourses 批量查询课程报名人数。
	EnrollNumByCourses(ctx context.Context, courseIds []int64) (map[int64]int64, error)
	// PurchaseInfo 查询课程购买信息。
	PurchaseInfo(ctx context.Context, courseId int64) (*types.PurchaseInfoVO, error)
	// PageOrderDetails 管理端分页查询订单明细。
	PageOrderDetails(ctx context.Context, q *OrderDetailQuery) (*types.Page, error)
	// GetOrderDetail 查询订单明细详情（管理端）。
	GetOrderDetail(ctx context.Context, id int64) (*types.OrderDetailAdminVO, error)
}

type orderDetailService struct {
	detailModel *model.OrderDetailModel
	orderModel  *model.OrderModel
	refundModel *model.RefundApplyModel
	userClient  UserClient
}

// NewOrderDetailService 创建订单明细业务服务。
func NewOrderDetailService(detailModel *model.OrderDetailModel, orderModel *model.OrderModel, refundModel *model.RefundApplyModel, userClient UserClient) OrderDetailService {
	return &orderDetailService{
		detailModel: detailModel,
		orderModel:  orderModel,
		refundModel: refundModel,
		userClient:  userClient,
	}
}

// IsCoursePurchased 用户是否已购买课程。
func (s *orderDetailService) IsCoursePurchased(ctx context.Context, userId, courseId int64) (bool, error) {
	if userId == 0 || courseId == 0 {
		return false, nil
	}
	detail, err := s.detailModel.FindPaidByUserCourse(ctx, userId, courseId)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, xerr.Wrap(err, xerr.CodeInternal, "查询购买记录失败")
	}
	return detail != nil, nil
}

// EnrollCourseAmounts 批量查询学员报名课程数。
func (s *orderDetailService) EnrollCourseAmounts(ctx context.Context, studentIds []int64) (map[int64]int64, error) {
	if len(studentIds) == 0 {
		return map[int64]int64{}, nil
	}
	return s.detailModel.CountEnrolledByUsers(ctx, studentIds)
}

// EnrollNumByCourses 批量查询课程报名人数。
func (s *orderDetailService) EnrollNumByCourses(ctx context.Context, courseIds []int64) (map[int64]int64, error) {
	if len(courseIds) == 0 {
		return map[int64]int64{}, nil
	}
	return s.detailModel.CountEnrolledByCourses(ctx, courseIds)
}

// PurchaseInfo 查询课程购买信息。
func (s *orderDetailService) PurchaseInfo(ctx context.Context, courseId int64) (*types.PurchaseInfoVO, error) {
	if courseId == 0 {
		return nil, xerr.BadRequestf("课程 id 不能为空")
	}
	info, err := s.detailModel.SumPurchaseInfo(ctx, courseId)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询课程购买信息失败")
	}
	return &types.PurchaseInfoVO{
		EnrollNum:     info.EnrollNum,
		RealPayAmount: info.RealPayAmount,
		RefundNum:     info.RefundNum,
	}, nil
}

// PageOrderDetails 管理端分页查询订单明细。
func (s *orderDetailService) PageOrderDetails(ctx context.Context, q *OrderDetailQuery) (*types.Page, error) {
	offset, limit := normalizePage(q.PageNo, q.PageSize)
	cond := &model.DetailPageCond{
		Id:           q.Id,
		Status:       q.Status,
		RefundStatus: q.RefundStatus,
		PayChannel:   q.PayChannel,
		StartTime:    q.StartTime,
		EndTime:      q.EndTime,
		Offset:       offset,
		Limit:        limit,
		IsAsc:        q.IsAsc,
	}
	if q.Mobile != "" {
		users, err := s.userClient.GetByIds(ctx, nil)
		_ = users
		_ = err
		// 用户服务暂不支持按手机号检索，过滤逻辑在接入后补充
	}

	rows, total, err := s.detailModel.ListPage(ctx, cond)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}

	userIDs := make(map[int64]struct{}, len(rows))
	for i := range rows {
		userIDs[rows[i].UserId] = struct{}{}
	}
	users := s.fetchUserMap(ctx, userIDs)

	list := make([]types.OrderDetailPageVO, 0, len(rows))
	for i := range rows {
		r := &rows[i]
		vo := types.OrderDetailPageVO{
			Id:               r.Id,
			OrderId:          r.OrderId,
			Name:             r.Name,
			Price:            r.Price,
			RealPayAmount:    r.RealPayAmount,
			PayChannel:       r.PayChannel,
			Status:           r.Status,
			StatusDesc:       detailStatusDesc(r.Status),
			RefundStatus:     formatNullInt(r.RefundStatus),
			RefundStatusDesc: refundStatusDesc(formatNullInt(r.RefundStatus)),
			CreateTime:       r.CreateTime.Format(time.RFC3339),
		}
		if u, ok := users[r.UserId]; ok {
			vo.Mobile = u.CellPhone
		}
		list = append(list, vo)
	}
	return &types.Page{List: list, Total: total, Pages: calcPages(total, limit)}, nil
}

// GetOrderDetail 查询订单明细详情（管理端）。
func (s *orderDetailService) GetOrderDetail(ctx context.Context, id int64) (*types.OrderDetailAdminVO, error) {
	detail, err := s.detailModel.FindById(ctx, id)
	if err == sql.ErrNoRows {
		return nil, xerr.NotFound("订单明细不存在")
	}
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询订单明细失败")
	}

	vo := &types.OrderDetailAdminVO{
		Id:             detail.Id,
		OrderId:        detail.OrderId,
		Name:           detail.Name,
		Price:          detail.Price,
		RealPayAmount:  detail.RealPayAmount,
		DiscountAmount: detail.DiscountAmount,
		PayChannel:     detail.PayChannel,
		Status:         detail.Status,
		RefundStatus:   formatNullInt(detail.RefundStatus),
		CanRefund:      canRefund(detail.Status, formatNullInt(detail.RefundStatus)),
		StudyValidTime: formatNullTime(detail.CourseExpireTime),
	}

	order, err := s.orderModel.FindById(ctx, detail.OrderId)
	if err == nil {
		vo.PayOrderNo = formatNullInt(order.PayOrderNo)
		vo.Message = order.Message
		vo.Nodes = orderProgressNodes(order)
	}

	apply, err := s.refundModel.FindByOrderDetail(ctx, detail.Id)
	if err == nil && apply != nil {
		vo.RefundApplyId = apply.Id
		vo.RefundReason = apply.RefundReason
		vo.RefundMessage = apply.Message
		vo.RefundChannel = apply.RefundChannel.String
		vo.RefundOrderNo = formatNullInt(apply.RefundOrderNo)
		vo.RefundProposerName = s.userName(ctx, apply.UserId)
		vo.FailedReason = apply.FailedReason.String
		vo.Remark = apply.Remark.String
	}

	if detail.UserId > 0 {
		users := s.fetchUserMap(ctx, map[int64]struct{}{detail.UserId: {}})
		if u, ok := users[detail.UserId]; ok {
			vo.StudentName = u.Name
			vo.Mobile = u.CellPhone
		}
	}
	return vo, nil
}

func (s *orderDetailService) fetchUserMap(ctx context.Context, ids map[int64]struct{}) map[int64]*UserInfo {
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

func (s *orderDetailService) userName(ctx context.Context, userId int64) string {
	if userId <= 0 {
		return ""
	}
	users := s.fetchUserMap(ctx, map[int64]struct{}{userId: {}})
	if u, ok := users[userId]; ok {
		return u.Name
	}
	return ""
}
