// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/course/internal/logic"
	"tjxt/api/course/internal/svc"
	"tjxt/api/course/internal/types"
	"tjxt/pkg/response"
)

func CourseSearchInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseSearchInfoLogic(r.Context(), svcCtx)
		resp, err := l.CourseSearchInfo(&req)
		response.Write(w, r, resp, err)
	}
}
