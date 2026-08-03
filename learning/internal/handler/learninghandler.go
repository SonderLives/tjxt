// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"learning/internal/logic"
	"learning/internal/svc"
	"learning/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func LearningHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLearningLogic(r.Context(), svcCtx)
		resp, _ := l.Learning(&req)
		httpx.OkJsonCtx(r.Context(), w, resp)
	}
}
