package model

// 优惠券状态
const (
	CouponStatusDraft  = "draft"  // 草稿（待发放）
	CouponStatusIssued = "issued" // 已发放
	CouponStatusPaused = "paused" // 暂停发放
	CouponStatusEnded  = "ended"  // 已结束
)

// 优惠券类型
const (
	DiscountTypeReduce      = "reduce"       // 满减
	DiscountTypeDiscount    = "discount"     // 折扣
	DiscountTypeNoThreshold = "no_threshold" // 无门槛
)

// 优惠券获取方式
const (
	ObtainWayReceive  = "receive"  // 手动领取
	ObtainWayExchange = "exchange" // 兑换码兑换
	ObtainWayAssign   = "assign"   // 后台发放
)

// 用户券状态
const (
	UserCouponStatusUnused  = "unused"  // 未使用
	UserCouponStatusUsed    = "used"    // 已使用
	UserCouponStatusExpired = "expired" // 已过期
)

// 兑换码状态
const (
	CouponCodeStatusUnused = "unused" // 未兑换
	CouponCodeStatusUsed   = "used"   // 已兑换
)
