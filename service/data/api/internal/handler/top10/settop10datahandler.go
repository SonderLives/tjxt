// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"net/http"

	"tjxt/service/data/api/internal/logic/top10"
	"tjxt/service/data/api/internal/svc"
	"tjxt/service/data/api/internal/types"
	"tjxt/common/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetTop10DataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Top10DataSetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := top10.NewSetTop10DataLogic(r.Context(), svcCtx)
		resp, err := l.SetTop10Data(&req)
		response.Write(w, r, resp, err)
	}
}
