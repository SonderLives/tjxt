package logic

import (
	"context"

	"tjxt/pkg/auth"
	"tjxt/pkg/response"
	"tjxt/apps/trade/api/internal/service"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 学员是否已购买课程 ============

type OrderDetailCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailCourseLogic {
	return &OrderDetailCourseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderDetailCourseLogic) OrderDetailCourse(req *types.OrderDetailCourseReq) (resp *result.R, err error) {
	userID, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	purchased, err := l.svcCtx.OrderDetailService.IsCoursePurchased(l.ctx, userID, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(purchased), nil
}

// ============ 批量查询学员报名课程数 ============

type EnrollCourseLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnrollCourseLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnrollCourseLogic {
	return &EnrollCourseLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *EnrollCourseLogic) EnrollCourse(req *types.EnrollCourseReq) (resp *result.R, err error) {
	studentIDs := parseIDs(req.StudentIds)
	if len(studentIDs) == 0 {
		return result.OkData(map[int64]int64{}), nil
	}
	amounts, err := l.svcCtx.OrderDetailService.EnrollCourseAmounts(l.ctx, studentIDs)
	if err != nil {
		return nil, err
	}
	return result.OkData(amounts), nil
}

// ============ 批量查询课程报名人数 ============

type EnrollNumLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewEnrollNumLogic(ctx context.Context, svcCtx *svc.ServiceContext) *EnrollNumLogic {
	return &EnrollNumLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *EnrollNumLogic) EnrollNum(req *types.EnrollNumReq) (resp *result.R, err error) {
	courseIDs := parseIDs(req.CourseIdList)
	if len(courseIDs) == 0 {
		return result.OkData(map[int64]int64{}), nil
	}
	nums, err := l.svcCtx.OrderDetailService.EnrollNumByCourses(l.ctx, courseIDs)
	if err != nil {
		return nil, err
	}
	return result.OkData(nums), nil
}

// ============ 订单明细分页（管理端） ============

type OrderDetailPageLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailPageLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailPageLogic {
	return &OrderDetailPageLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderDetailPageLogic) OrderDetailPage(req *types.OrderDetailPageReq) (resp *result.R, err error) {
	q := &service.OrderDetailQuery{
		Id:           req.Id,
		Mobile:       req.Mobile,
		Status:       req.Status,
		RefundStatus: req.RefundStatus,
		PayChannel:   req.PayChannel,
		PageNo:       req.PageNo,
		PageSize:     req.PageSize,
		IsAsc:        req.IsAsc,
	}
	if t, ok := parseTime(req.OrderStartTime); ok {
		q.StartTime = t
	}
	if t, ok := parseTime(req.OrderEndTime); ok {
		q.EndTime = t
	}
	page, err := l.svcCtx.OrderDetailService.PageOrderDetails(l.ctx, q)
	if err != nil {
		return nil, err
	}
	return result.OkData(page), nil
}

// ============ 课程购买信息 ============

type PurchaseInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPurchaseInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PurchaseInfoLogic {
	return &PurchaseInfoLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *PurchaseInfoLogic) PurchaseInfo(req *types.PurchaseInfoReq) (resp *result.R, err error) {
	info, err := l.svcCtx.OrderDetailService.PurchaseInfo(l.ctx, req.CourseId)
	if err != nil {
		return nil, err
	}
	return result.OkData(info), nil
}

// ============ 订单明细详情（管理端） ============

type OrderDetailGetLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewOrderDetailGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *OrderDetailGetLogic {
	return &OrderDetailGetLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *OrderDetailGetLogic) OrderDetailGet(req *types.OrderIdReq) (resp *result.R, err error) {
	detail, err := l.svcCtx.OrderDetailService.GetOrderDetail(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(detail), nil
}
