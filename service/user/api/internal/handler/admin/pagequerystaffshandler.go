// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package admin

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/user/api/internal/logic/admin"
	"tjxt/service/user/api/internal/svc"
	"tjxt/service/user/api/internal/types"
	"tjxt/common/response"
)

func PageQueryStaffsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := admin.NewPageQueryStaffsLogic(r.Context(), svcCtx)
		resp, err := l.PageQueryStaffs(&req)
		response.Write(w, r, resp, err)
	}
}
