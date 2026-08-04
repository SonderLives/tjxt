package handler

import "net/http"

// LoginIP 从请求头中提取客户端 IP 地址
// 优先级：X-Forwarded-For 第一个 IP > X-Real-IP > RemoteAddr
func LoginIP(r *http.Request) string {
	// 1. X-Forwarded-For: 可能包含多个 IP，取第一个
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		// XFF 格式：client, proxy1, proxy2；取第一个
		idx := 0
		for i, c := range xff {
			if c == ',' {
				idx = i
				break
			}
		}
		return xff[:idx]
	}
	// 2. X-Real-IP
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	// 3. RemoteAddr
	return r.RemoteAddr
}