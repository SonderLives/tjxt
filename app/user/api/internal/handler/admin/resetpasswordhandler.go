// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/user/api/internal/logic/admin"
	"tjxt/service/user/api/internal/svc"
	"tjxt/service/user/api/internal/types"
	"tjxt/common/response"
)

func ResetPasswordHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewResetPasswordLogic(r.Context(), svcCtx)
		err := l.ResetPassword(&req)
		response.Write(w, r, nil, err)
	}
}
