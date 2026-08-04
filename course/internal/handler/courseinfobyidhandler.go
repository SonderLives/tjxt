// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"strconv"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

func CourseInfoByIdHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := pathvar.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var req types.CourseInfoQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseInfoByIdLogic(r.Context(), svcCtx)
		resp, err := l.CourseInfoById(id, &req)
		writeResult(w, r, resp, err)
	}
}
