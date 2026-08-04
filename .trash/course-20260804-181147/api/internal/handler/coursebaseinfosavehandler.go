// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/course/api/internal/logic"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CourseBaseInfoSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseBaseInfoSaveDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseBaseInfoSaveLogic(r.Context(), svcCtx)
		resp, err := l.CourseBaseInfoSave(&req)
		writeResult(w, r, resp, err)
	}
}
