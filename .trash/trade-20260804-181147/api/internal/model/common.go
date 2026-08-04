package model

import (
	"database/sql"
	"time"
)

// nullInt64 将 sql.NullInt64 转为可写参数。
func nullInt64(n sql.NullInt64) any {
	if !n.Valid {
		return nil
	}
	return n.Int64
}

// nullString 将 sql.NullString 转为可写参数。
func nullString(n sql.NullString) any {
	if !n.Valid {
		return nil
	}
	return n.String
}

// nullTime 将 sql.NullTime 转为可写参数。
func nullTime(n sql.NullTime) any {
	if !n.Valid {
		return nil
	}
	return n.Time
}

// NowPtr 返回当前时间指针，便于写入 time.Time 字段。
func NowPtr() *time.Time {
	now := time.Now()
	return &now
}
