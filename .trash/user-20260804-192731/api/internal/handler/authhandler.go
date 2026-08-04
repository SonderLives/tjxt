package handler

import (
	"net/http"

	"tjxt/apps/user/api/internal/logic"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// bearerTokenFrom 从请求 Authorization 头提取 Bearer token。
func bearerTokenFrom(r *http.Request) string {
	const prefix = "Bearer "
	h := r.Header.Get("Authorization")
	if len(h) > len(prefix) && h[:len(prefix)] == prefix {
		return h[len(prefix):]
	}
	return h
}

func AccountLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginFormDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewAccountLoginLogic(r.Context(), svcCtx)
		resp, err := l.AccountLogin(&req)
		writeResult(w, r, resp, err)
	}
}

func AdminLoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginFormDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewAdminLoginLogic(r.Context(), svcCtx)
		resp, err := l.AdminLogin(&req)
		writeResult(w, r, resp, err)
	}
}

func AccountRefreshHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAccountRefreshLogic(r.Context(), svcCtx, bearerTokenFrom(r))
		resp, err := l.AccountRefresh()
		writeResult(w, r, resp, err)
	}
}

func AccountLogoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAccountLogoutLogic(r.Context(), svcCtx)
		resp, err := l.AccountLogout()
		writeResult(w, r, resp, err)
	}
}
