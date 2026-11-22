// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/api/pay/internal/logic"
	"tjxt/api/pay/internal/svc"
	"tjxt/api/pay/internal/types"
	"tjxt/pkg/response"
)

func QueryRefundResultHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QueryRefundResultReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewQueryRefundResultLogic(r.Context(), svcCtx)
		resp, err := l.QueryRefundResult(&req)
		response.Write(w, r, resp, err)
	}
}
