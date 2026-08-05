package logic

import (
	"strconv"

	"tjxt/apps/promotion/api/internal/types"
	pb "tjxt/apps/promotion/rpc/promotion"
)

// toPageRequest 将 API 层分页请求转换为 RPC 分页结构。
func toPageRequest(in types.PageRequest) *pb.PageRequest {
	return &pb.PageRequest{
		PageNo:   in.PageNo,
		PageSize: in.PageSize,
		IsAsc:    in.IsAsc,
		SortBy:   in.SortBy,
	}
}

// toCouponDetailVO pb 优惠券详情 -> API 视图对象。
func toCouponDetailVO(in *pb.CouponDetailVO) *types.CouponDetailVO {
	if in == nil {
		return nil
	}
	return &types.CouponDetailVO{
		Id:                in.Id,
		Name:              in.Name,
		DiscountType:      in.DiscountType,
		DiscountValue:     in.DiscountValue,
		MaxDiscountAmount: in.MaxDiscountAmount,
		ThresholdAmount:   in.ThresholdAmount,
		ObtainWay:         in.ObtainWay,
		Specific:          in.Specific,
		IssueBeginTime:    in.IssueBeginTime,
		IssueEndTime:      in.IssueEndTime,
		TermBeginTime:     in.TermBeginTime,
		TermDays:          in.TermDays,
		TermEndTime:       in.TermEndTime,
		TotalNum:          in.TotalNum,
		UserLimit:         in.UserLimit,
		Status:            in.Status,
	}
}

// toCouponPageVO pb 优惠券分页项 -> API 视图对象。
func toCouponPageVO(in *pb.CouponPageVO) *types.CouponPageVO {
	if in == nil {
		return nil
	}
	return &types.CouponPageVO{
		Id:                in.Id,
		Name:              in.Name,
		DiscountType:      in.DiscountType,
		DiscountValue:     in.DiscountValue,
		MaxDiscountAmount: in.MaxDiscountAmount,
		ThresholdAmount:   in.ThresholdAmount,
		Specific:          in.Specific,
		ObtainWay:         in.ObtainWay,
		TotalNum:          in.TotalNum,
		IssueNum:          in.IssueNum,
		UsedNum:           in.UsedNum,
		Status:            in.Status,
		TermDays:          in.TermDays,
		TermBeginTime:     in.TermBeginTime,
		TermEndTime:       in.TermEndTime,
		IssueBeginTime:    in.IssueBeginTime,
		IssueEndTime:      in.IssueEndTime,
		CreateTime:        in.CreateTime,
	}
}

// toCouponVO pb 优惠券视图 -> API 视图对象。
func toCouponVO(in *pb.CouponVO) *types.CouponVO {
	if in == nil {
		return nil
	}
	return &types.CouponVO{
		Id:                in.Id,
		Name:              in.Name,
		DiscountType:      in.DiscountType,
		DiscountValue:     in.DiscountValue,
		MaxDiscountAmount: in.MaxDiscountAmount,
		ThresholdAmount:   in.ThresholdAmount,
		Specific:          in.Specific,
		TermDays:          in.TermDays,
		TermEndTime:       in.TermEndTime,
		Available:         in.Available,
		Received:          in.Received,
	}
}

// toExchangeCodeVO pb 兑换码 -> API 视图对象。
func toExchangeCodeVO(in *pb.ExchangeCodeVO) *types.ExchangeCodeVO {
	if in == nil {
		return nil
	}
	return &types.ExchangeCodeVO{Id: in.Id, Code: in.Code}
}

// toCouponDiscountVO pb 折扣计算结果 -> API 视图对象。
// DiscountDetail 的键为课程 id（int64），API 层以字符串呈现。
func toCouponDiscountVO(in *pb.CouponDiscountDTO) *types.CouponDiscountVO {
	if in == nil {
		return nil
	}
	detail := make(map[string]int64, len(in.DiscountDetail))
	for k, v := range in.DiscountDetail {
		detail[strconv.FormatInt(k, 10)] = v
	}
	return &types.CouponDiscountVO{
		DiscountAmount: in.DiscountAmount,
		DiscountDetail: detail,
		Ids:            in.Ids,
		Rules:          in.Rules,
	}
}
