// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/service/promotion/api/internal/logic"
	"tjxt/service/promotion/api/internal/svc"
	"tjxt/common/response"
)

func CouponListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCouponListLogic(r.Context(), svcCtx)
		resp, err := l.CouponList()
		response.Write(w, r, resp, err)
	}
}
