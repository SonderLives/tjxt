// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"net/http"

	"tjxt/api/data/internal/logic/today"
	"tjxt/api/data/internal/svc"
	"tjxt/api/data/internal/types"
	"tjxt/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetTodayDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TodayDataSetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := today.NewSetTodayDataLogic(r.Context(), svcCtx)
		resp, err := l.SetTodayData(&req)
		response.Write(w, r, resp, err)
	}
}
