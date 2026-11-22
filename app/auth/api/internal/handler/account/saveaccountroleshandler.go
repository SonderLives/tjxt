// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/auth/api/internal/logic/account"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"
	"tjxt/common/response"
)

func SaveAccountRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AccountRoleReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := account.NewSaveAccountRolesLogic(r.Context(), svcCtx)
		resp, err := l.SaveAccountRoles(&req)
		response.Write(w, r, resp, err)
	}
}
