// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/promotion/api/internal/logic"
	"tjxt/apps/promotion/api/internal/svc"
	"tjxt/apps/promotion/api/internal/types"
	result "tjxt/pkg/response"
)

func CouponGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCouponGetLogic(r.Context(), svcCtx)
		resp, err := l.CouponGet(&req)
		result.Write(w, r, resp, err)
	}
}
