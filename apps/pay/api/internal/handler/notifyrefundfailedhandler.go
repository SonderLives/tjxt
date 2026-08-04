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

func NotifyRefundFailedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NotifyRefundFailedReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewNotifyRefundFailedLogic(r.Context(), svcCtx)
		resp, err := l.NotifyRefundFailed(&req)
		result.Write(w, r, resp, err)
	}
}
