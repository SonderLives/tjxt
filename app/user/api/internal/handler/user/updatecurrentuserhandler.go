// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/user/api/internal/logic/user"
	"tjxt/app/user/api/internal/svc"
	"tjxt/app/user/api/internal/types"
	"tjxt/pkg/response"
)

func UpdateCurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserFormReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := user.NewUpdateCurrentUserLogic(r.Context(), svcCtx)
		err := l.UpdateCurrentUser(&req)
		response.Write(w, r, nil, err)
	}
}
