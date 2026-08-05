package logic

import (
	"sort"
	"strconv"
	"strings"
	"time"

	"tjxt/apps/promotion/rpc/internal/model"
	"tjxt/apps/promotion/rpc/pb"
)

// maxCombineCoupons 参与组合运算的最大券数量。
// 组合方案数为 2^n，超出部分按单券优惠力度截断，避免组合爆炸拖垮接口。
const maxCombineCoupons = 12

// maxSolutions 最多返回的方案数。
const maxSolutions = 30

// availableCoupon 一张可参与计算的用户券（用户券 id + 券规则）。
type availableCoupon struct {
	userCouponId int64
	coupon       *model.Coupon
}

// usableUserCoupons 组装当前用户可用的券：未使用、未过期、且券规则仍存在。
func usableUserCoupons(ucs []*model.UserCoupon, coupons map[int64]*model.Coupon, now time.Time) []availableCoupon {
	result := make([]availableCoupon, 0, len(ucs))
	for _, uc := range ucs {
		if uc.Status != model.UserCouponStatusUnused {
			continue
		}
		if uc.ExpireTime.Valid && now.After(uc.ExpireTime.Time) {
			continue
		}
		c, ok := coupons[uc.CouponId]
		if !ok || c.Deleted == 1 {
			continue
		}
		// 有效期未开始的券暂不可用
		if c.TermBeginTime.Valid && now.Before(c.TermBeginTime.Time) {
			continue
		}
		result = append(result, availableCoupon{userCouponId: uc.Id, coupon: c})
	}
	return result
}

// matchedCourses 返回券适用范围内的课程。
func matchedCourses(c *model.Coupon, courses []*pb.OrderCourseDTO) []*pb.OrderCourseDTO {
	if c.Specific == 0 {
		return courses
	}
	matched := make([]*pb.OrderCourseDTO, 0, len(courses))
	for _, course := range courses {
		if matchScope(c, course.CateId) {
			matched = append(matched, course)
		}
	}
	return matched
}

// buildSolution 按给定券组合计算优惠方案。
// 券按顺序作用于「剩余金额」，避免同一笔钱被重复打折；
// 组合中任一张券无法生效（未达门槛或无适用课程）时返回 false，即该组合非法。
func buildSolution(combo []availableCoupon, courses []*pb.OrderCourseDTO) (*pb.CouponDiscountDTO, bool) {
	if len(combo) == 0 {
		return nil, false
	}

	remaining := make(map[int64]int64, len(courses))
	for _, course := range courses {
		remaining[course.Id] += course.Price
	}

	var (
		total  int64
		detail = make(map[int64]int64, len(courses))
		ids    = make([]int64, 0, len(combo))
		rules  = make([]string, 0, len(combo))
	)

	for _, ac := range combo {
		matched := matchedCourses(ac.coupon, courses)
		if len(matched) == 0 {
			return nil, false
		}

		var subtotal int64
		for _, course := range matched {
			subtotal += remaining[course.Id]
		}

		discount := calcDiscount(ac.coupon, subtotal)
		if discount <= 0 {
			return nil, false
		}

		// 按课程剩余金额比例分摊优惠，最后一门课兜底吸收取整误差
		var allocated int64
		for i, course := range matched {
			var share int64
			if i == len(matched)-1 {
				share = discount - allocated
			} else {
				share = discount * remaining[course.Id] / subtotal
			}
			if share > remaining[course.Id] {
				share = remaining[course.Id]
			}
			remaining[course.Id] -= share
			detail[course.Id] += share
			allocated += share
		}

		total += allocated
		ids = append(ids, ac.userCouponId)
		rules = append(rules, couponRule(ac.coupon))
	}

	if total <= 0 {
		return nil, false
	}
	return &pb.CouponDiscountDTO{
		DiscountAmount: total,
		DiscountDetail: detail,
		Ids:            ids,
		Rules:          rules,
	}, true
}

// trimCoupons 券数量超限时，按单券优惠力度保留最有价值的若干张。
func trimCoupons(available []availableCoupon, courses []*pb.OrderCourseDTO) []availableCoupon {
	if len(available) <= maxCombineCoupons {
		return available
	}
	type scored struct {
		ac    availableCoupon
		value int64
	}
	list := make([]scored, 0, len(available))
	for _, ac := range available {
		s, ok := buildSolution([]availableCoupon{ac}, courses)
		var v int64
		if ok {
			v = s.DiscountAmount
		}
		list = append(list, scored{ac: ac, value: v})
	}
	sort.SliceStable(list, func(i, j int) bool { return list[i].value > list[j].value })

	result := make([]availableCoupon, 0, maxCombineCoupons)
	for i := 0; i < maxCombineCoupons; i++ {
		result = append(result, list[i].ac)
	}
	return result
}

// calcSolutions 穷举券组合，返回按优惠金额降序排列的可用方案。
func calcSolutions(available []availableCoupon, courses []*pb.OrderCourseDTO) []*pb.CouponDiscountDTO {
	if len(available) == 0 || len(courses) == 0 {
		return nil
	}
	available = trimCoupons(available, courses)

	var (
		solutions = make([]*pb.CouponDiscountDTO, 0, 16)
		seen      = make(map[string]struct{})
		n         = len(available)
	)

	// 位掩码穷举所有非空子集
	for mask := 1; mask < (1 << n); mask++ {
		combo := make([]availableCoupon, 0, n)
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				combo = append(combo, available[i])
			}
		}
		solution, ok := buildSolution(combo, courses)
		if !ok {
			continue
		}
		key := solutionKey(solution.Ids)
		if _, dup := seen[key]; dup {
			continue
		}
		seen[key] = struct{}{}
		solutions = append(solutions, solution)
	}

	// 优惠多的排前面；金额相同时用券少的更优
	sort.SliceStable(solutions, func(i, j int) bool {
		if solutions[i].DiscountAmount != solutions[j].DiscountAmount {
			return solutions[i].DiscountAmount > solutions[j].DiscountAmount
		}
		return len(solutions[i].Ids) < len(solutions[j].Ids)
	})

	if len(solutions) > maxSolutions {
		solutions = solutions[:maxSolutions]
	}
	return solutions
}

// solutionKey 以券 id 集合作为方案去重键。
func solutionKey(ids []int64) string {
	sorted := append([]int64{}, ids...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	var sb strings.Builder
	for _, id := range sorted {
		sb.WriteString(strconv.FormatInt(id, 10))
		sb.WriteByte(',')
	}
	return sb.String()
}
