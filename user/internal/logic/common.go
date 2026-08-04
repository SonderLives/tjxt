package logic

import (
	"context"
	"strings"

	"common/auth"
	"common/xerr"
)

// currentUserID 从上下文提取当前登录用户 id。
func currentUserID(ctx context.Context) (int64, error) {
	return auth.UserIdFromCtx(ctx)
}

// requireStaffOrAdmin 管理端操作权限校验。
func requireStaffOrAdmin(ctx context.Context) error {
	if !auth.IsStaffOrAdmin(ctx) {
		return xerr.Forbidden("无权操作")
	}
	return nil
}

// bearerToken 从 Authorization 头提取 Bearer token。
func bearerToken(header string) string {
	const prefix = "Bearer "
	h := strings.TrimSpace(header)
	if strings.HasPrefix(h, prefix) {
		return strings.TrimSpace(h[len(prefix):])
	}
	return h
}
