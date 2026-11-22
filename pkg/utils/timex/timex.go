// Package timex 提供全仓库统一的时间字符串格式化/解析工具。
//
// 此前 7 个 RPC 服务（course/exam/learning/message/pay/promotion/trade 等）
// 在各自的 logic/common.go 里重复定义了 "2006-01-02 15:04:05" 布局与
// formatTime/formatNullTime/parseDateTime 等同义函数；本包将其收敛为
// 一份实现，供服务重构时统一替换。
package timex

import (
	"database/sql"
	"time"
)

// Layout 全仓库统一的时间字符串布局：yyyy-MM-dd HH:mm:ss。
// 与各 RPC 服务 logic/common.go 中原有 timeLayout 常量一致。
const Layout = "2006-01-02 15:04:05"

// Format 将时间格式化为 yyyy-MM-dd HH:mm:ss 字符串。
//
// 零值时间返回空串（而非 "0001-01-01 00:00:00"）——保持与各服务原有
// formatTime 的语义一致：零值在 VO/展示层表示"暂无时间"。
func Format(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format(Layout)
}

// Parse 按 Layout 解析时间字符串，使用本地时区（与原 parseDateTime 一致，
// 避免前端传入的壁钟时间被当成 UTC）。空串或非法格式返回 error。
func Parse(s string) (time.Time, error) {
	return time.ParseInLocation(Layout, s, time.Local)
}

// FormatNullTime 格式化 sql.NullTime：无效值返回空串，有效值同 Format。
// 与各服务原有 formatNullTime 语义一致。
func FormatNullTime(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}
	return Format(t.Time)
}

// ParseNullTime 解析时间字符串为 sql.NullTime：
// 空串返回无效 NullTime 与 nil error（等价于写库 NULL，与原 parseDateTime 一致）；
// 非空串解析失败返回错误；解析成功返回有效 NullTime。
func ParseNullTime(s string) (sql.NullTime, error) {
	if s == "" {
		return sql.NullTime{}, nil
	}
	t, err := Parse(s)
	if err != nil {
		return sql.NullTime{}, err
	}
	return sql.NullTime{Time: t, Valid: true}, nil
}
