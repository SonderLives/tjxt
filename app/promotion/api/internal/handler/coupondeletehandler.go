// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/promotion/api/internal/logic"
	"tjxt/app/promotion/api/internal/svc"
	"tjxt/app/promotion/api/internal/types"
	"tjxt/pkg/response"
)

func CouponDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCouponDeleteLogic(r.Context(), svcCtx)
		err := l.CouponDelete(&req)
		response.Write(w, r, nil, err)
	}
}
