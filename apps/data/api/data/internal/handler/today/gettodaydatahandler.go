// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package today

import (
	"net/http"

	"tjxt/apps/data/api/data/internal/logic/today"
	"tjxt/apps/data/api/data/internal/svc"
	result "tjxt/pkg/response"
)

func GetTodayDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := today.NewGetTodayDataLogic(r.Context(), svcCtx)
		resp, err := l.GetTodayData()
		result.Write(w, r, resp, err)
	}
}
