// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/exam/internal/logic/question"
	"tjxt/api/exam/internal/svc"
	"tjxt/api/exam/internal/types"
	"tjxt/pkg/response"
)

func ListQuestionsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question.NewListQuestionsLogic(r.Context(), svcCtx)
		resp, err := l.ListQuestions(&req)
		response.Write(w, r, resp, err)
	}
}
