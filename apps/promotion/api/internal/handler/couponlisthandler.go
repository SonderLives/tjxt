// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/promotion/api/internal/logic"
	"tjxt/apps/promotion/api/internal/svc"
	result "tjxt/pkg/response"
)

func CouponListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCouponListLogic(r.Context(), svcCtx)
		resp, err := l.CouponList()
		result.Write(w, r, resp, err)
	}
}
