// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/search/api/internal/logic/interests"
	"tjxt/service/search/api/internal/svc"
	"tjxt/service/search/api/internal/types"
	"tjxt/common/response"
)

func SaveInterestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SaveInterestsReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := interests.NewSaveInterestsLogic(r.Context(), svcCtx)
		resp, err := l.SaveInterests(&req)
		response.Write(w, r, resp, err)
	}
}
