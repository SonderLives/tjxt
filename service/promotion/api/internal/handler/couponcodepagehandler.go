// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/promotion/internal/logic"
	"tjxt/api/promotion/internal/svc"
	"tjxt/api/promotion/internal/types"
	"tjxt/pkg/response"
)

func CouponCodePageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponCodePageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCouponCodePageLogic(r.Context(), svcCtx)
		resp, err := l.CouponCodePage(&req)
		response.Write(w, r, resp, err)
	}
}
