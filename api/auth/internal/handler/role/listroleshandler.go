// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/auth/internal/logic/role"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"
	"tjxt/pkg/response"
)

func ListRolesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewListRolesLogic(r.Context(), svcCtx)
		resp, err := l.ListRoles(&req)
		response.Write(w, r, resp, err)
	}
}
