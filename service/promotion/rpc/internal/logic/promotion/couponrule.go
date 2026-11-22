package promotionlogic

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"tjxt/rpc/promotion/internal/model"
)

// parseScopes 解析 scopes JSON（形如 [1,2,3]）为 id 列表。
func parseScopes(s sql.NullString) []int64 {
	if !s.Valid || strings.TrimSpace(s.String) == "" {
		return nil
	}
	var ids []int64
	if err := json.Unmarshal([]byte(s.String), &ids); err != nil {
		return nil
	}
	return ids
}

// matchScope 判断课程分类是否落在优惠券的适用范围内。
// specific=0 表示全场通用。
func matchScope(c *model.Coupon, cateId int64) bool {
	if c.Specific == 0 {
		return true
	}
	for _, id := range parseScopes(c.Scopes) {
		if id == cateId {
			return true
		}
	}
	return false
}

// yuan 分转元，用于生成面向用户的规则文案。
func yuan(cent int64) string {
	return strings.TrimSuffix(strings.TrimRight(fmt.Sprintf("%.2f", float64(cent)/100), "0"), ".")
}

// couponRule 生成优惠券规则描述，如「满100元减20元」。
func couponRule(c *model.Coupon) string {
	switch c.DiscountType {
	case model.DiscountTypeNoThreshold:
		return fmt.Sprintf("无门槛立减%s元", yuan(c.DiscountValue))
	case model.DiscountTypeDiscount:
		text := fmt.Sprintf("满%s元打%s折", yuan(c.ThresholdAmount), yuan(c.DiscountValue*10))
		if c.MaxDiscountAmount > 0 {
			text += fmt.Sprintf("，最多减%s元", yuan(c.MaxDiscountAmount))
		}
		return text
	default: // reduce
		return fmt.Sprintf("满%s元减%s元", yuan(c.ThresholdAmount), yuan(c.DiscountValue))
	}
}

// calcDiscount 按券规则计算订单可优惠金额（单位：分）。
// 未达门槛返回 0；折扣券受 max_discount_amount 封顶；优惠不超过订单总额。
func calcDiscount(c *model.Coupon, totalAmount int64) int64 {
	if totalAmount <= 0 {
		return 0
	}
	if c.DiscountType != model.DiscountTypeNoThreshold && totalAmount < c.ThresholdAmount {
		return 0
	}

	var discount int64
	switch c.DiscountType {
	case model.DiscountTypeDiscount:
		// discountValue 为折扣百分比，如 80 表示 8 折，优惠额 = 总额 *(100-80)/100
		if c.DiscountValue <= 0 || c.DiscountValue >= 100 {
			return 0
		}
		discount = totalAmount * (100 - c.DiscountValue) / 100
		if c.MaxDiscountAmount > 0 && discount > c.MaxDiscountAmount {
			discount = c.MaxDiscountAmount
		}
	default: // reduce / no_threshold
		discount = c.DiscountValue
	}

	if discount > totalAmount {
		discount = totalAmount
	}
	return discount
}

// couponReceivable 判断优惠券当前是否可被领取：发放中、在发放期内、且未领完。
func couponReceivable(c *model.Coupon, now time.Time) bool {
	if c.Status != model.CouponStatusIssued {
		return false
	}
	if c.IssueBeginTime.Valid && now.Before(c.IssueBeginTime.Time) {
		return false
	}
	if c.IssueEndTime.Valid && now.After(c.IssueEndTime.Time) {
		return false
	}
	if c.TotalNum > 0 && c.IssueNum >= c.TotalNum {
		return false
	}
	return true
}

// userCouponExpireTime 计算用户券的过期时间：优先使用绝对有效期，其次按领取日 +termDays。
func userCouponExpireTime(c *model.Coupon, now time.Time) sql.NullTime {
	if c.TermEndTime.Valid {
		return c.TermEndTime
	}
	if c.TermDays > 0 {
		end := now.AddDate(0, 0, int(c.TermDays))
		// 有效期末尾对齐到当天 23:59:59，避免按秒卡点体验差
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())
		return sql.NullTime{Time: end, Valid: true}
	}
	return sql.NullTime{}
}
