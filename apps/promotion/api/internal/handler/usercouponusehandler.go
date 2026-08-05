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

func UserCouponUseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CouponIdsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUserCouponUseLogic(r.Context(), svcCtx)
		err := l.UserCouponUse(&req)
		result.Write(w, r, nil, err)
	}
}
