// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/exam/api/internal/logic/questionbiz"
	"tjxt/service/exam/api/internal/svc"
	"tjxt/service/exam/api/internal/types"
	"tjxt/common/response"
)

func RemoveQuestionBizHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuestionBizReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := questionbiz.NewRemoveQuestionBizLogic(r.Context(), svcCtx)
		resp, err := l.RemoveQuestionBiz(&req)
		response.Write(w, r, resp, err)
	}
}
