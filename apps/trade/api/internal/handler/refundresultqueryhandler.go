// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/trade/api/internal/logic"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"
	result "tjxt/pkg/response"
)

func RefundResultQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BizRefundOrderIdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewRefundResultQueryLogic(r.Context(), svcCtx)
		resp, err := l.RefundResultQuery(&req)
		result.Write(w, r, resp, err)
	}
}
