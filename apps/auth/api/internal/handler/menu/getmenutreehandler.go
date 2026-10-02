// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	"tjxt/apps/auth/api/internal/logic/menu"
	"tjxt/apps/auth/api/internal/svc"
	result "tjxt/pkg/response"
)

func GetMenuTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := menu.NewGetMenuTreeLogic(r.Context(), svcCtx)
		resp, err := l.GetMenuTree()
		result.Write(w, r, resp, err)
	}
}
