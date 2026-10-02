// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"net/http"

	"tjxt/apps/search/api/internal/logic/interests"
	"tjxt/apps/search/api/internal/svc"
	result "tjxt/pkg/response"
)

func GetInterestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := interests.NewGetInterestsLogic(r.Context(), svcCtx)
		resp, err := l.GetInterests()
		result.Write(w, r, resp, err)
	}
}
