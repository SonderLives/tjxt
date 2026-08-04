// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/trade/api/internal/logic"
	"tjxt/apps/trade/api/internal/svc"
	result "tjxt/pkg/response"
)

func CartListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCartListLogic(r.Context(), svcCtx)
		resp, err := l.CartList()
		result.Write(w, r, resp, err)
	}
}
