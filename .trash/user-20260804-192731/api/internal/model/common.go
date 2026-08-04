package model

import (
	"context"
	"database/sql"
)

// 用户类型常量
const (
	UserTypeStaff   int64 = 1 // 员工
	UserTypeStudent int64 = 2 // 学员
	UserTypeTeacher int64 = 3 // 老师
)

// 账户状态
const (
	UserStatusDisabled int64 = 0 // 禁用
	UserStatusNormal   int64 = 1 // 正常
)

// nullString 可空字符串
type nullString = sql.NullString

// nullInt64 可空整型
type nullInt64 = sql.NullInt64

// execer 兼容 sqlx.SqlConn 与 sqlx.Session 的写执行器。
type execer interface {
	ExecCtx(ctx context.Context, query string, args ...any) (sql.Result, error)
}
