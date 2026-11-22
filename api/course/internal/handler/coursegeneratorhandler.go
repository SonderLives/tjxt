// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/api/course/internal/logic"
	"tjxt/api/course/internal/svc"
	"tjxt/pkg/response"
)

func CourseGeneratorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCourseGeneratorLogic(r.Context(), svcCtx)
		resp, err := l.CourseGenerator()
		response.Write(w, r, resp, err)
	}
}
