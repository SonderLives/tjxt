package model

import "strings"

// inPlaceholders 依据 id 列表生成 `?,?,?` 占位符及对应参数，供 IN 查询使用，
// 避免把 id 直接拼进 SQL 造成注入风险。调用方需自行保证 ids 非空。
func inPlaceholders(ids []int64) (string, []any) {
	holders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	return holders, args
}

// pairValuePlaceholders 生成 n 组 `(?,?)` 的批量插入占位符。
func pairValuePlaceholders(n int) string {
	return strings.TrimSuffix(strings.Repeat("(?,?),", n), ",")
}

// dedupeIds 去重并剔除非正数 id，保持原有顺序。
func dedupeIds(ids []int64) []int64 {
	if len(ids) == 0 {
		return nil
	}
	seen := make(map[int64]struct{}, len(ids))
	out := make([]int64, 0, len(ids))
	for _, id := range ids {
		if id <= 0 {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}
