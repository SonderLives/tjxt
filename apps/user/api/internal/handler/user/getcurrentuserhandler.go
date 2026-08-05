// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package user

import (
	"net/http"

	"tjxt/apps/user/api/internal/logic/user"
	"tjxt/apps/user/api/internal/svc"
	result "tjxt/pkg/response"
)

func GetCurrentUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := user.NewGetCurrentUserLogic(r.Context(), svcCtx)
		resp, err := l.GetCurrentUser()
		result.Write(w, r, resp, err)
	}
}
