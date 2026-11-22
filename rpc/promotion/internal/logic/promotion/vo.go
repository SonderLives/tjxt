package promotionlogic

import (
	"database/sql"
	"encoding/json"
	"strings"

	"tjxt/pkg/errorx"
	"tjxt/pkg/timex"
	"tjxt/rpc/promotion/internal/model"
	"tjxt/rpc/promotion/pb"
)

// toCouponDetailVO 券管理端详情 VO。
func toCouponDetailVO(c *model.Coupon) *pb.CouponDetailVO {
	scopeIds := parseScopes(c.Scopes)
	scopes := make([]*pb.CouponScopeVO, 0, len(scopeIds))
	for _, id := range scopeIds {
		// 分类名称由 course 服务维护，此处仅回传 id，避免跨服务强依赖
		scopes = append(scopes, &pb.CouponScopeVO{Id: id})
	}
	return &pb.CouponDetailVO{
		Id:                c.Id,
		Name:              c.Name,
		DiscountType:      c.DiscountType,
		DiscountValue:     c.DiscountValue,
		MaxDiscountAmount: c.MaxDiscountAmount,
		ThresholdAmount:   c.ThresholdAmount,
		ObtainWay:         c.ObtainWay,
		Specific:          c.Specific == 1,
		Scopes:            scopes,
		IssueBeginTime:    timex.FormatNullTime(c.IssueBeginTime),
		IssueEndTime:      timex.FormatNullTime(c.IssueEndTime),
		TermBeginTime:     timex.FormatNullTime(c.TermBeginTime),
		TermDays:          c.TermDays,
		TermEndTime:       timex.FormatNullTime(c.TermEndTime),
		Status:            c.Status,
		TotalNum:          c.TotalNum,
		UserLimit:         c.UserLimit,
	}
}

func toCouponPageVO(c *model.Coupon) *pb.CouponPageVO {
	return &pb.CouponPageVO{
		Id:                c.Id,
		Name:              c.Name,
		DiscountType:      c.DiscountType,
		DiscountValue:     c.DiscountValue,
		MaxDiscountAmount: c.MaxDiscountAmount,
		ThresholdAmount:   c.ThresholdAmount,
		Specific:          c.Specific == 1,
		ObtainWay:         c.ObtainWay,
		TotalNum:          c.TotalNum,
		IssueNum:          c.IssueNum,
		UsedNum:           c.UsedNum,
		Status:            c.Status,
		TermDays:          c.TermDays,
		TermBeginTime:     timex.FormatNullTime(c.TermBeginTime),
		TermEndTime:       timex.FormatNullTime(c.TermEndTime),
		IssueBeginTime:    timex.FormatNullTime(c.IssueBeginTime),
		IssueEndTime:      timex.FormatNullTime(c.IssueEndTime),
		CreateTime:        timex.Format(c.CreateTime),
	}
}

// toCouponVO 面向 C 端的券信息。available 表示当前是否可领，received 表示当前用户是否已领过。
func toCouponVO(c *model.Coupon, available, received bool) *pb.CouponVO {
	return &pb.CouponVO{
		Id:                c.Id,
		Name:              c.Name,
		DiscountType:      c.DiscountType,
		DiscountValue:     c.DiscountValue,
		MaxDiscountAmount: c.MaxDiscountAmount,
		ThresholdAmount:   c.ThresholdAmount,
		Specific:          c.Specific == 1,
		TermDays:          c.TermDays,
		TermEndTime:       timex.FormatNullTime(c.TermEndTime),
		Available:         available,
		Received:          received,
	}
}

// buildCoupon 把表单转换为 model，新建的券固定为草稿状态。
func buildCoupon(in *pb.CouponFormDTO, operator int64) *model.Coupon {
	specific := int64(0)
	scopes := sql.NullString{}
	if in.Specific {
		specific = 1
		if len(in.Scopes) > 0 {
			if buf, err := json.Marshal(in.Scopes); err == nil {
				scopes = sql.NullString{String: string(buf), Valid: true}
			}
		}
	}
	return &model.Coupon{
		Id:                in.Id,
		Name:              strings.TrimSpace(in.Name),
		DiscountType:      in.DiscountType,
		DiscountValue:     in.DiscountValue,
		MaxDiscountAmount: in.MaxDiscountAmount,
		ThresholdAmount:   in.ThresholdAmount,
		ObtainWay:         in.ObtainWay,
		Specific:          specific,
		Scopes:            scopes,
		TotalNum:          in.TotalNum,
		UserLimit:         in.UserLimit,
		Status:            model.CouponStatusDraft,
		Creater:           operator,
		Updater:           operator,
	}
}

// validateCouponForm 校验优惠券表单的业务约束。
func validateCouponForm(in *pb.CouponFormDTO) error {
	if strings.TrimSpace(in.Name) == "" {
		return errorx.BadRequestf("优惠券名称不能为空")
	}
	switch in.DiscountType {
	case model.DiscountTypeReduce:
		if in.ThresholdAmount <= 0 {
			return errorx.BadRequestf("满减券必须设置使用门槛")
		}
		if in.DiscountValue <= 0 || in.DiscountValue >= in.ThresholdAmount {
			return errorx.BadRequestf("满减金额必须大于 0 且小于使用门槛")
		}
	case model.DiscountTypeDiscount:
		if in.DiscountValue <= 0 || in.DiscountValue >= 100 {
			return errorx.BadRequestf("折扣值必须在 1~99 之间")
		}
	case model.DiscountTypeNoThreshold:
		if in.DiscountValue <= 0 {
			return errorx.BadRequestf("无门槛券优惠金额必须大于 0")
		}
	default:
		return errorx.BadRequestf("优惠券类型非法：%s", in.DiscountType)
	}
	switch in.ObtainWay {
	case model.ObtainWayReceive, model.ObtainWayExchange, model.ObtainWayAssign:
	default:
		return errorx.BadRequestf("获取方式非法：%s", in.ObtainWay)
	}
	if in.TotalNum < 0 || in.UserLimit < 0 {
		return errorx.BadRequestf("发行总量与每人限领数量不能为负数")
	}
	if in.Specific && len(in.Scopes) == 0 {
		return errorx.BadRequestf("限定适用范围时必须指定范围")
	}
	return nil
}
