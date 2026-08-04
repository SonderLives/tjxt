package logic

import (
	"strconv"
	"strings"
	"time"

	"common/xerr"
)

// errBadRequest 构造参数错误。
func errBadRequest(msg string) error {
	return xerr.BadRequestf("%s", msg)
}

// parseIDs 将逗号分隔的 id 字符串解析为 int64 切片，非法元素忽略。
func parseIDs(s string) []int64 {
	var ids []int64
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if v, err := strconv.ParseInt(part, 10, 64); err == nil && v > 0 {
			ids = append(ids, v)
		}
	}
	return ids
}

// parseTime 解析 RFC3339 时间串，空串返回零值与 false。
func parseTime(s string) (time.Time, bool) {
	if strings.TrimSpace(s) == "" {
		return time.Time{}, false
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}, false
	}
	return t, true
}
