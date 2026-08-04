// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"
	"strconv"

	"tjxt/apps/course/api/internal/logic"
	"tjxt/apps/course/api/internal/svc"

	"github.com/zeromicro/go-zero/rest/pathvar"
)

func CourseCatasIndexListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := pathvar.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		l := logic.NewCourseCatasIndexListLogic(r.Context(), svcCtx)
		resp, err := l.CourseCatasIndexList(id)
		writeResult(w, r, resp, err)
	}
}
