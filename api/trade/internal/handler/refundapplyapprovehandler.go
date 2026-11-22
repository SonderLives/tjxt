// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/trade/internal/logic"
	"tjxt/api/trade/internal/svc"
	"tjxt/api/trade/internal/types"
	"tjxt/pkg/response"
)

func RefundApplyApproveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApproveFormReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRefundApplyApproveLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyApprove(&req)
		response.Write(w, r, resp, err)
	}
}
