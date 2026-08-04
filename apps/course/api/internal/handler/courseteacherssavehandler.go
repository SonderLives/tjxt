// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/course/api/internal/logic"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	result "tjxt/pkg/response"
)

func CourseTeachersSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTeacherSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseTeachersSaveLogic(r.Context(), svcCtx)
		resp, err := l.CourseTeachersSave(&req)
		result.Write(w, r, resp, err)
	}
}
