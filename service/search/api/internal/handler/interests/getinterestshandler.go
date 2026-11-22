// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package interests

import (
	"net/http"

	"tjxt/service/search/api/internal/logic/interests"
	"tjxt/service/search/api/internal/svc"
	"tjxt/common/response"
)

func GetInterestsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := interests.NewGetInterestsLogic(r.Context(), svcCtx)
		resp, err := l.GetInterests()
		response.Write(w, r, resp, err)
	}
}
