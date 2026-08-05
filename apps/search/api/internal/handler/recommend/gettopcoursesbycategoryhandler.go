// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package recommend

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/search/api/internal/logic/recommend"
	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"
)

func GetTopCoursesByCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := recommend.NewGetTopCoursesByCategoryLogic(r.Context(), svcCtx)
		resp, err := l.GetTopCoursesByCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
