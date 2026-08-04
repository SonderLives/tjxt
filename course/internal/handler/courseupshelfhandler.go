// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"strconv"

	"course/internal/logic"
	"course/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func CourseUpShelfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := pathvar.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		l := logic.NewCourseUpShelfLogic(r.Context(), svcCtx)
		resp, err := l.UpShelf(id)
		writeResult(w, r, resp, err)
	}
}
