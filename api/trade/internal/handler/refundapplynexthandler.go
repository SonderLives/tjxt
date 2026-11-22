// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/trade/api/internal/logic"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/pkg/response"
)

func RefundApplyNextHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewRefundApplyNextLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyNext()
		response.Write(w, r, resp, err)
	}
}
