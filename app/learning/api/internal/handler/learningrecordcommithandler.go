// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/learning/api/internal/logic"
	"tjxt/service/learning/api/internal/svc"
	"tjxt/service/learning/api/internal/types"
	"tjxt/common/response"
)

func LearningRecordCommitHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LearningRecordCommitReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewLearningRecordCommitLogic(r.Context(), svcCtx)
		resp, err := l.LearningRecordCommit(&req)
		response.Write(w, r, resp, err)
	}
}
