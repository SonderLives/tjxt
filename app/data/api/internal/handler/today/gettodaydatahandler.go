// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"net/http"

	"tjxt/app/data/api/internal/logic/today"
	"tjxt/app/data/api/internal/svc"
	"tjxt/pkg/response"
)

func GetTodayDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := today.NewGetTodayDataLogic(r.Context(), svcCtx)
		resp, err := l.GetTodayData()
		response.Write(w, r, resp, err)
	}
}
