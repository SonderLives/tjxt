// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/user/internal/logic/admin"
	"tjxt/api/user/internal/svc"
	"tjxt/api/user/internal/types"
	"tjxt/pkg/response"
)

func UpdateUserByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserUpdateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewUpdateUserByIdLogic(r.Context(), svcCtx)
		err := l.UpdateUserById(&req)
		response.Write(w, r, nil, err)
	}
}
