package logic

import (
	"strconv"
	"strings"
	"time"

	"tjxt/apps/course/api/internal/types"
	"tjxt/apps/course/rpc/pb"
)

// formatTime 时间格式化（空值返回空串）。
func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// parseIds 将逗号分隔的 id 字符串解析为 int64 切片（忽略非法片段）。
func parseIds(s string) []int64 {
	var ids []int64
	for _, p := range strings.Split(s, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if v, err := strconv.ParseInt(p, 10, 64); err == nil {
			ids = append(ids, v)
		}
	}
	return ids
}

// toCategoryVO pb.CategoryInfo -> API CategoryVO
func toCategoryVO(c *pb.CategoryInfo) types.CategoryVO {
	return types.CategoryVO{
		Id:               c.Id,
		Name:             c.Name,
		ParentId:         c.ParentId,
		Level:            int64(c.Level),
		Index:            int64(c.Priority),
		Status:           int64(c.Status),
		StatusDesc:       c.StatusDesc,
		CourseNum:        c.CourseNum,
		ThirdCategoryNum: c.ThirdCategoryNum,
		CreateTime:       c.CreateTime,
		UpdateTime:       c.UpdateTime,
	}
}

// toSimpleCategoryVO pb.CategoryNode -> API SimpleCategoryVO（递归）
func toSimpleCategoryVO(n *pb.CategoryNode) *types.SimpleCategoryVO {
	if n == nil {
		return nil
	}
	vo := &types.SimpleCategoryVO{
		Id:       n.Id,
		Name:     n.Name,
		ParentId: n.ParentId,
		Level:    int64(n.Level),
	}
	for _, child := range n.Children {
		vo.Children = append(vo.Children, *toSimpleCategoryVO(child))
	}
	return vo
}
