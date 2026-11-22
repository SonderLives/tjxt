// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package question

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/exam/api/internal/logic/question"
	"tjxt/service/exam/api/internal/svc"
	"tjxt/service/exam/api/internal/types"
	"tjxt/common/response"
)

func SaveQuestionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := question.NewSaveQuestionLogic(r.Context(), svcCtx)
		resp, err := l.SaveQuestion(&req)
		response.Write(w, r, resp, err)
	}
}
