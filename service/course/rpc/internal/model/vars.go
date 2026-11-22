package model

import (
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var ErrNotFound = sqlx.ErrNotFound

// insertColumns 统计生成模型列名片段（如 courseRowsExpectAutoSet）的列数。
func insertColumns(columns string) int {
	return strings.Count(columns, ",") + 1
}

// rowPlaceholders 生成 n 列的单行插入占位符 "(?, ?, ...)"，
// 与生成模型列名片段配套使用，列数自动同步。
func rowPlaceholders(n int) string {
	return "(" + strings.TrimSuffix(strings.Repeat("?,", n), ",") + ")"
}
