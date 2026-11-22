// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/promotion/api/internal/logic"
	"tjxt/service/promotion/api/internal/svc"
	"tjxt/service/promotion/api/internal/types"
	"tjxt/common/response"
)

func CouponIssueHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponIssueReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCouponIssueLogic(r.Context(), svcCtx)
		err := l.CouponIssue(&req)
		response.Write(w, r, nil, err)
	}
}
