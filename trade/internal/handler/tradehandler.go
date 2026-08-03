// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"go-zero/internal/response"
	"trade/internal/logic"
	"trade/internal/svc"
	"trade/internal/types"
)

func TradeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewTradeLogic(r.Context(), svcCtx)
		resp, err := l.Trade(&req)
		response.Response(w, resp, err)
	}
}
