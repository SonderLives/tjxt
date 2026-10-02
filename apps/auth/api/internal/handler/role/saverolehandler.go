// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package role

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/auth/api/internal/logic/role"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"
	result "tjxt/pkg/response"
)

func SaveRoleHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RoleSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := role.NewSaveRoleLogic(r.Context(), svcCtx)
		resp, err := l.SaveRole(&req)
		result.Write(w, r, resp, err)
	}
}
