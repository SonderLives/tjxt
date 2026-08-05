// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/user/api/internal/logic/admin"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"
	result "tjxt/pkg/response"
)

func AddUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewAddUserLogic(r.Context(), svcCtx)
		resp, err := l.AddUser(&req)
		result.Write(w, r, resp, err)
	}
}
