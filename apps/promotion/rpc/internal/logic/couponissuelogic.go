package logic

import (
	"context"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/internal/svc"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CouponIssueLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCouponIssueLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CouponIssueLogic {
	return &CouponIssueLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// CouponIssue 发放优惠券：写入发放期与有效期并置为发放中。
// 兑换码类型的券在首次发放时按发行总量生成兑换码。
func (l *CouponIssueLogic) CouponIssue(in *pb.CouponIssueFormDTO) (*pb.Empty, error) {
	if in.Id <= 0 {
		return nil, xerr.BadRequestf("优惠券 id 非法")
	}

	data, err := l.svcCtx.CouponModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if isNotFound(err) {
			return nil, xerr.NotFound("优惠券不存在")
		}
		l.Errorf("find coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "发放优惠券失败")
	}
	if data.Deleted == 1 {
		return nil, xerr.NotFound("优惠券不存在")
	}
	if data.Status == model.CouponStatusEnded {
		return nil, xerr.Conflict("优惠券已结束，无法再次发放")
	}

	issueBegin, err := parseTime(in.IssueBeginTime)
	if err != nil {
		return nil, err
	}
	issueEnd, err := parseTime(in.IssueEndTime)
	if err != nil {
		return nil, err
	}
	termBegin, err := parseTime(in.TermBeginTime)
	if err != nil {
		return nil, err
	}
	termEnd, err := parseTime(in.TermEndTime)
	if err != nil {
		return nil, err
	}

	if issueBegin.Valid && issueEnd.Valid && issueEnd.Time.Before(issueBegin.Time) {
		return nil, xerr.BadRequestf("发放结束时间不能早于开始时间")
	}
	if termBegin.Valid && termEnd.Valid && termEnd.Time.Before(termBegin.Time) {
		return nil, xerr.BadRequestf("使用结束时间不能早于开始时间")
	}
	// 有效期二选一：绝对区间 或 相对天数
	if in.TermDays <= 0 && !termEnd.Valid {
		return nil, xerr.BadRequestf("请设置有效期天数或使用结束时间")
	}

	firstIssue := data.Status == model.CouponStatusDraft

	data.IssueBeginTime = issueBegin
	data.IssueEndTime = issueEnd
	data.TermBeginTime = termBegin
	data.TermEndTime = termEnd
	data.TermDays = in.TermDays
	data.Status = model.CouponStatusIssued
	// 未指定发放开始时间视为立即发放
	if !data.IssueBeginTime.Valid {
		data.IssueBeginTime = sqlNow()
	}

	if err := l.svcCtx.CouponModel.Update(l.ctx, data); err != nil {
		l.Errorf("issue coupon failed, id=%d, err=%v", in.Id, err)
		return nil, xerr.Wrap(err, xerr.CodeInternal, "发放优惠券失败")
	}

	// 兑换码券首次发放时批量生成兑换码
	if firstIssue && data.ObtainWay == model.ObtainWayExchange && data.TotalNum > 0 {
		codes, err := generateCodes(data.TotalNum)
		if err != nil {
			return nil, err
		}
		if err := l.svcCtx.CouponCodeModel.BatchInsert(l.ctx, data.Id, codes, data.Updater); err != nil {
			l.Errorf("batch insert coupon code failed, couponId=%d, err=%v", data.Id, err)
			return nil, xerr.Wrap(err, xerr.CodeInternal, "生成兑换码失败")
		}
	}

	return &pb.Empty{}, nil
}
