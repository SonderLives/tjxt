// Package auth 提供从请求上下文提取登录态与 JWT 签发/校验工具。
//
// JWT 载荷约定（与 learning 服务一致）：
//   - userId: 用户主键（必填，数值型）
//   - role:   角色，USER/STUDENT/TEACHER/ADMIN 等，缺省视为 USER
//
// 该约定由 auth 服务签发时遵循；各业务服务仅消费，不自行签发访问令牌。
package auth

import (
	"context"
	"fmt"
	"strconv"

	"common/xerr"
)

// Context 中的 claim 键名（go-zero rest.WithJwt 注入 ctx 时以 JWT 载荷字段名为 key，类型为 string，
// 因此这里必须用字符串字面量做 key，才能与中间件注入值匹配）。
const (
	keyUserID = "userId"
	keyRole   = "role"
)

// Role 常量
const (
	RoleUser    = "USER"
	RoleStudent = "STUDENT"
	RoleTeacher = "TEACHER"
	RoleStaff   = "STAFF"
	RoleAdmin   = "ADMIN"
)

// UserIdFromCtx 从上下文提取当前登录用户 id。
func UserIdFromCtx(ctx context.Context) (int64, error) {
	v := ctx.Value(keyUserID)
	if v == nil {
		return 0, xerr.Unauthorized("")
	}
	userID, err := strconv.ParseInt(fmt.Sprint(v), 10, 64)
	if err != nil || userID <= 0 {
		return 0, xerr.Unauthorized("")
	}
	return userID, nil
}

// RoleFromCtx 从上下文提取角色，缺失时返回 USER。
func RoleFromCtx(ctx context.Context) string {
	v := ctx.Value(keyRole)
	if v == nil {
		return RoleUser
	}
	if s, ok := v.(string); ok && s != "" {
		return s
	}
	return fmt.Sprint(v)
}

// IsAdmin 判断当前用户是否管理员。
func IsAdmin(ctx context.Context) bool {
	return RoleFromCtx(ctx) == RoleAdmin
}

// IsStaffOrAdmin 判断当前用户是否为员工或管理员（管理端操作）。
func IsStaffOrAdmin(ctx context.Context) bool {
	switch RoleFromCtx(ctx) {
	case RoleStaff, RoleAdmin:
		return true
	}
	return false
}
