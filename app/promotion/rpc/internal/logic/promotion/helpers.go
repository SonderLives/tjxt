package promotionlogic

import (
	"database/sql"
	"strings"
	"time"

	"tjxt/common/errorx"
	"tjxt/common/timex"
)

// sqlNow 返回当前时间的 sql.NullTime。
func sqlNow() sql.NullTime {
	return sql.NullTime{Time: time.Now(), Valid: true}
}

// parseTime 解析前端传入的时间字符串，兼容 "2006-01-02 15:04:05" 与 "2006-01-02" 等格式。
// 空串返回无效 NullTime（等价于写库 NULL）；无法识别的格式返回 BadRequest。
func parseTime(s string) (sql.NullTime, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return sql.NullTime{}, nil
	}
	if t, err := timex.Parse(s); err == nil {
		return sql.NullTime{Time: t, Valid: true}, nil
	}
	for _, layout := range []string{"2006-01-02T15:04:05", time.RFC3339, "2006-01-02"} {
		if t, err := time.ParseInLocation(layout, s, time.Local); err == nil {
			return sql.NullTime{Time: t, Valid: true}, nil
		}
	}
	return sql.NullTime{}, errorx.BadRequestf("时间格式非法：%s", s)
}
