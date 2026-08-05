// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"net/http"

	"tjxt/apps/data/api/data/internal/logic/today"
	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	result "tjxt/pkg/response"

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
		result.Write(w, r, resp, err)
	}
}
