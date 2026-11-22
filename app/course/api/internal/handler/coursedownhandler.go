// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/course/api/internal/logic"
	"tjxt/app/course/api/internal/svc"
	"tjxt/app/course/api/internal/types"
	"tjxt/pkg/response"
)

func CourseDownHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseIdListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseDownLogic(r.Context(), svcCtx)
		resp, err := l.CourseDown(&req)
		response.Write(w, r, resp, err)
	}
}
