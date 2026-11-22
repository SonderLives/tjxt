// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/auth/internal/logic/privilege"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"
	"tjxt/pkg/response"
)

func SavePrivilegeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PrivilegeSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := privilege.NewSavePrivilegeLogic(r.Context(), svcCtx)
		resp, err := l.SavePrivilege(&req)
		response.Write(w, r, resp, err)
	}
}
