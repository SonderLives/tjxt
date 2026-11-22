package tradelogic

import (
	"context"
	"time"

	courseclient "tjxt/rpc/course/client/course"
	"tjxt/rpc/trade/internal/svc"
)

// ===================== 工具函数 =====================

func now() time.Time {
	return time.Now()
}

// ===================== 课程信息回填 =====================

// fetchCourseMap 批量查询课程简况，返回 course_id -> 课程信息
func fetchCourseMap(ctx context.Context, svcCtx *svc.ServiceContext, ids []int64) map[int64]*courseclient.CourseSimpleInfoItem {
	res := make(map[int64]*courseclient.CourseSimpleInfoItem)
	if len(ids) == 0 {
		return res
	}
	reply, err := svcCtx.CourseRpc.CourseSimpleInfoList(ctx, &courseclient.CourseSimpleInfoQueryRequest{Ids: ids})
	if err != nil {
		return res
	}
	for _, it := range reply.Items {
		res[it.Id] = it
	}
	return res
}

func courseName(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) string {
	if it, ok := m[id]; ok {
		return it.Name
	}
	return ""
}

func courseCover(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) string {
	if it, ok := m[id]; ok {
		return it.CoverUrl
	}
	return ""
}

func coursePrice(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) int64 {
	if it, ok := m[id]; ok {
		return it.Price
	}
	return 0
}

// courseThirdCate 返回课程三级分类 id（优惠券适用范围匹配用）。
func courseThirdCate(m map[int64]*courseclient.CourseSimpleInfoItem, id int64) int64 {
	if it, ok := m[id]; ok {
		return it.ThirdCateId
	}
	return 0
}
