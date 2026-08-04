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

func CourseCatasHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := pathvar.Vars(r)
		idStr := vars["id"]
		id, _ := strconv.ParseInt(idStr, 10, 64)

		var req types.CourseCatasQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseCatasLogic(r.Context(), svcCtx)
		resp, err := l.Catas(id, &req)
		writeResult(w, r, resp, err)
	}
}
