// Package service 课程微服务业务层。
// 业务规则与 Java 参考实现（tjxt）保持一致：目录树拼接、上下架、媒资/题目/老师保存等。
package service

import (
	"database/sql"
	"time"
)

// parseTimeOrNow 解析 "2006-01-02 15:04:05" 时间字符串，空值返回当前时间。
func parseTimeOrNow(s string) time.Time {
	t, err := time.ParseInLocation(timeFormat, s, time.Local)
	if err != nil {
		return time.Now()
	}
	return t
}

// formatNullTime 格式化可空时间字段，NULL 返回空字符串。
func formatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return t.Time.Format(timeFormat)
}

// distinctIDs 去重并返回正整数 id 列表。
func distinctIDs(ids []int64) []int64 {
	set := make(map[int64]bool, len(ids))
	result := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 || set[id] {
			continue
		}
		set[id] = true
		result = append(result, id)
	}
	return result
}

// containsInt64 判断列表中是否包含目标值。
func containsInt64(list []int64, target int64) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}

// 课程分类状态（与 tjxt CommonStatus 一致）
const (
	CategoryDisable = int64(0) // 禁用
	CategoryEnable  = int64(1) // 启用
)

// 课程状态（与 tjxt CourseStatus 一致）
const (
	CourseStatusNoUpShelf = int64(1) // 待上架
	CourseStatusShelf     = int64(2) // 已上架
	CourseStatusDownShelf = int64(3) // 下架
	CourseStatusFinished  = int64(4) // 已完结
)

// 课程编辑步骤（与 tjxt CourseConstants.CourseStep 一致）
const (
	CourseStepBaseInfo  = int64(1) // 基本信息
	CourseStepCatalogue = int64(2) // 目录
	CourseStepMedia     = int64(3) // 视频
	CourseStepSubject   = int64(4) // 题目
	CourseStepTeacher   = int64(5) // 老师
)

// 目录类型（与 tjxt CourseConstants.CataType 一致）
const (
	CataTypeChapter = int64(1) // 章
	CataTypeSection = int64(2) // 节
	CataTypePractice = int64(3) // 练习或测试
)

// 一级分类根节点 id
const CategoryRoot = int64(0)

// bool 值转 0/1
func boolToInt(v bool) int64 {
	if v {
		return 1
	}
	return 0
}

// 0/1 转 bool
func intToBool(v int64) bool {
	return v == 1
}

// 时间格式化为字符串（数据库 loc=Local，返回 "2006-01-02 15:04:05"）
const timeFormat = "2006-01-02 15:04:05"
