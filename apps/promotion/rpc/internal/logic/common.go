package logic

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
	"time"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/pb"
	"tjxt/pkg/xerr"
)

const timeLayout = "2006-01-02 15:04:05"

// codeAlphabet 兑换码字符集，剔除易混淆的 0/O/1/I。
const codeAlphabet = "23456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// codeLength 兑换码长度。
const codeLength = 12

// generateCodes 生成 n 个互不相同的兑换码；DB 上 uk_code 唯一索引兜底防重。
func generateCodes(n int64) ([]string, error) {
	if n <= 0 {
		return nil, nil
	}
	seen := make(map[string]struct{}, n)
	codes := make([]string, 0, n)
	max := big.NewInt(int64(len(codeAlphabet)))

	for int64(len(codes)) < n {
		var sb strings.Builder
		sb.Grow(codeLength)
		for i := 0; i < codeLength; i++ {
			idx, err := rand.Int(rand.Reader, max)
			if err != nil {
				return nil, xerr.Wrap(err, xerr.CodeInternal, "生成兑换码失败")
			}
			sb.WriteByte(codeAlphabet[idx.Int64()])
		}
		code := sb.String()
		if _, ok := seen[code]; ok {
			continue
		}
		seen[code] = struct{}{}
		codes = append(codes, code)
	}
	return codes, nil
}

// ============================ 时间转换 ============================

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(timeLayout)
}

func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timeLayout)
}

func formatNullString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	return s.String
}

// sqlNow 返回当前时间的 sql.NullTime。
func sqlNow() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

// parseTime 解析前端传入的时间字符串，兼容 "2006-01-02 15:04:05" 与 "2006-01-02"。
func parseTime(s string) (sql.NullTime, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullTime{}, nil
	}
	for _, layout := range []string{timeLayout, "2006-01-02T15:04:05", time.RFC3339, "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return sql.NullTime{Time: t, Valid: true}, nil
		}
	}
	return sql.NullTime{}, xerr.BadRequestf("时间格式非法：%s", s)
}

// ============================ 错误判定 ============================

func isNotFound(err error) bool {
	return err == sql.ErrNoRows || err == model.ErrNotFound
}

// ============================ 适用范围 ============================

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

// ============================ 优惠券规则 ============================

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

// ============================ model -> pb ============================

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
		IssueBeginTime:    formatNullTime(c.IssueBeginTime),
		IssueEndTime:      formatNullTime(c.IssueEndTime),
		TermBeginTime:     formatNullTime(c.TermBeginTime),
		TermDays:          c.TermDays,
		TermEndTime:       formatNullTime(c.TermEndTime),
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
		TermBeginTime:     formatNullTime(c.TermBeginTime),
		TermEndTime:       formatNullTime(c.TermEndTime),
		IssueBeginTime:    formatNullTime(c.IssueBeginTime),
		IssueEndTime:      formatNullTime(c.IssueEndTime),
		CreateTime:        formatTime(c.CreateTime),
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
		TermEndTime:       formatNullTime(c.TermEndTime),
		Available:         available,
		Received:          received,
	}
}

// ============================ pb -> model ============================

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
		return xerr.BadRequestf("优惠券名称不能为空")
	}
	switch in.DiscountType {
	case model.DiscountTypeReduce:
		if in.ThresholdAmount <= 0 {
			return xerr.BadRequestf("满减券必须设置使用门槛")
		}
		if in.DiscountValue <= 0 || in.DiscountValue >= in.ThresholdAmount {
			return xerr.BadRequestf("满减金额必须大于 0 且小于使用门槛")
		}
	case model.DiscountTypeDiscount:
		if in.DiscountValue <= 0 || in.DiscountValue >= 100 {
			return xerr.BadRequestf("折扣值必须在 1~99 之间")
		}
	case model.DiscountTypeNoThreshold:
		if in.DiscountValue <= 0 {
			return xerr.BadRequestf("无门槛券优惠金额必须大于 0")
		}
	default:
		return xerr.BadRequestf("优惠券类型非法：%s", in.DiscountType)
	}
	switch in.ObtainWay {
	case model.ObtainWayReceive, model.ObtainWayExchange, model.ObtainWayAssign:
	default:
		return xerr.BadRequestf("获取方式非法：%s", in.ObtainWay)
	}
	if in.TotalNum < 0 || in.UserLimit < 0 {
		return xerr.BadRequestf("发行总量与每人限领数量不能为负数")
	}
	if in.Specific && len(in.Scopes) == 0 {
		return xerr.BadRequestf("限定适用范围时必须指定范围")
	}
	return nil
}
