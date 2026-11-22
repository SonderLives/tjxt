// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/auth/api/internal/logic/role"
	"tjxt/app/auth/api/internal/svc"
	"tjxt/app/auth/api/internal/types"
	"tjxt/pkg/response"
)

func DeleteRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewDeleteRoleLogic(r.Context(), svcCtx)
		resp, err := l.DeleteRole(&req)
		response.Write(w, r, resp, err)
	}
}
