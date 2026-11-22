// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/learning/internal/logic"
	"tjxt/api/learning/internal/svc"
	"tjxt/api/learning/internal/types"
	"tjxt/pkg/response"
)

func LearningRecordsByCourseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LessonRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLearningRecordsByCourseLogic(r.Context(), svcCtx)
		resp, err := l.LearningRecordsByCourse(&req)
		response.Write(w, r, resp, err)
	}
}
