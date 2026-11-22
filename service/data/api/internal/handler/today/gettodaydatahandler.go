// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"net/http"

	"tjxt/service/data/api/internal/logic/today"
	"tjxt/service/data/api/internal/svc"
	"tjxt/common/response"
)

func GetTodayDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := today.NewGetTodayDataLogic(r.Context(), svcCtx)
		resp, err := l.GetTodayData()
		response.Write(w, r, resp, err)
	}
}
