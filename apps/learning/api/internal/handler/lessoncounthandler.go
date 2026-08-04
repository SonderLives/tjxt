// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/learning/api/internal/logic"
	"tjxt/apps/learning/api/internal/svc"
	"tjxt/apps/learning/api/internal/types"
	result "tjxt/pkg/response"
)

func LessonCountHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LessonRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLessonCountLogic(r.Context(), svcCtx)
		resp, err := l.LessonCount(&req)
		result.Write(w, r, resp, err)
	}
}
