// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"net/http"

	"tjxt/api/search/internal/logic/interests"
	"tjxt/api/search/internal/svc"
	"tjxt/pkg/response"
)

func GetInterestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := interests.NewGetInterestsLogic(r.Context(), svcCtx)
		resp, err := l.GetInterests()
		response.Write(w, r, resp, err)
	}
}
