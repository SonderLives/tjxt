// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/learning/api/internal/logic"
	"tjxt/apps/learning/api/internal/svc"
	result "tjxt/pkg/response"
)

func LearningNowHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewLearningNowLogic(r.Context(), svcCtx)
		resp, err := l.LearningNow()
		result.Write(w, r, resp, err)
	}
}
