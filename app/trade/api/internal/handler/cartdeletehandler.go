// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/trade/api/internal/logic"
	"tjxt/app/trade/api/internal/svc"
	"tjxt/app/trade/api/internal/types"
	"tjxt/pkg/response"
)

func CartDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CartIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCartDeleteLogic(r.Context(), svcCtx)
		resp, err := l.CartDelete(&req)
		response.Write(w, r, resp, err)
	}
}
