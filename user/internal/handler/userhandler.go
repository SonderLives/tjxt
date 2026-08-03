// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"go-zero/internal/response"
	"user/internal/logic"
	"user/internal/svc"
	"user/internal/types"
)

func UserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserLogic(r.Context(), svcCtx)
		resp, err := l.User(&req)
		response.Response(w, resp, err)
	}
}
