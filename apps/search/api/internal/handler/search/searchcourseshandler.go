// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package search

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/search/api/internal/logic/search"
	"tjxt/apps/search/api/internal/svc"
	"tjxt/apps/search/api/internal/types"
)

func SearchCoursesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseSearchReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := search.NewSearchCoursesLogic(r.Context(), svcCtx)
		resp, err := l.SearchCourses(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
