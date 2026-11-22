// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/course/api/internal/logic"
	"tjxt/service/course/api/internal/svc"
	"tjxt/service/course/api/internal/types"
	"tjxt/common/response"
)

func CourseSubjectsGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseSubjectsGetLogic(r.Context(), svcCtx)
		resp, err := l.CourseSubjectsGet(&req)
		response.Write(w, r, resp, err)
	}
}
