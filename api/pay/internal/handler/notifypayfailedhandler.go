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

func NotifyPayFailedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyPayFailedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewNotifyPayFailedLogic(r.Context(), svcCtx)
		resp, err := l.NotifyPayFailed(&req)
		response.Write(w, r, resp, err)
	}
}
