// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/auth/internal/logic/menu"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"
	"tjxt/pkg/response"
)

func SaveMenuHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MenuSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menu.NewSaveMenuLogic(r.Context(), svcCtx)
		resp, err := l.SaveMenu(&req)
		response.Write(w, r, resp, err)
	}
}
