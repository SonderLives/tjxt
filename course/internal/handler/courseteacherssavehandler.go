// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CourseTeachersSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTeacherSaveDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseTeachersSaveLogic(r.Context(), svcCtx)
		resp, err := l.CourseTeachersSave(&req)
		writeResult(w, r, resp, err)
	}
}
