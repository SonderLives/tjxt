// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/exam/api/internal/logic/question"
	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"
)

func DeleteQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question.NewDeleteQuestionLogic(r.Context(), svcCtx)
		resp, err := l.DeleteQuestion(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
