// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package assign

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/auth/api/internal/logic/assign"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"
	"tjxt/common/response"
)

func GetRolePrivilegesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := assign.NewGetRolePrivilegesLogic(r.Context(), svcCtx)
		resp, err := l.GetRolePrivileges(&req)
		response.Write(w, r, resp, err)
	}
}
