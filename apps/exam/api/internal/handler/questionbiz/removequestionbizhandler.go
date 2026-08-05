// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package questionbiz

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/exam/api/internal/logic/questionbiz"
	"tjxt/apps/exam/api/internal/svc"
	"tjxt/apps/exam/api/internal/types"
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
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
