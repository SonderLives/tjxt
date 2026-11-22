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

func OrderDetailEnrollCourseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnrollCourseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewOrderDetailEnrollCourseLogic(r.Context(), svcCtx)
		resp, err := l.OrderDetailEnrollCourse(&req)
		response.Write(w, r, resp, err)
	}
}
