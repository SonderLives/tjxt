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

func LessonPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLessonPageLogic(r.Context(), svcCtx)
		resp, err := l.LessonPage(&req)
		result.Write(w, r, resp, err)
	}
}
