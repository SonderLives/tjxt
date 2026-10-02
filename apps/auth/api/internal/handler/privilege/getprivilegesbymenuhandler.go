// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/auth/api/internal/logic/privilege"
	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"
	result "tjxt/pkg/response"
)

func GetPrivilegesByMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := privilege.NewGetPrivilegesByMenuLogic(r.Context(), svcCtx)
		resp, err := l.GetPrivilegesByMenu(&req)
		result.Write(w, r, resp, err)
	}
}
