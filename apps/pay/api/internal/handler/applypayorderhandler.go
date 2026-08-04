// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	result "tjxt/pkg/response"
	"tjxt/apps/pay/api/internal/logic"
	"tjxt/apps/pay/api/internal/svc"
	"tjxt/apps/pay/api/internal/types"
)

func ApplyPayOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApplyPayOrderReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewApplyPayOrderLogic(r.Context(), svcCtx)
		resp, err := l.ApplyPayOrder(&req)
		result.Write(w, r, resp, err)
	}
}
