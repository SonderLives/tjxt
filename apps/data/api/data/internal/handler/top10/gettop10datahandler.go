// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package top10

import (
	"net/http"

	"tjxt/apps/data/api/data/internal/logic/top10"
	"tjxt/apps/data/api/data/internal/svc"
	result "tjxt/pkg/response"
)

func GetTop10DataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := top10.NewGetTop10DataLogic(r.Context(), svcCtx)
		resp, err := l.GetTop10Data()
		result.Write(w, r, resp, err)
	}
}
