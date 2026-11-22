// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	"tjxt/service/auth/api/internal/logic/menu"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/common/response"
)

func GetMenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetMenuTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetMenuTree()
		response.Write(w, r, resp, err)
	}
}
