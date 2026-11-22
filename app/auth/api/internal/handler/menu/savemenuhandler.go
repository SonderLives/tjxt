// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package menu

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/auth/api/internal/logic/menu"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"
	"tjxt/common/response"
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
