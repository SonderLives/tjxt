// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package smsplatform

import (
	"net/http"

	"tjxt/apps/message/api/internal/logic/smsplatform"
	"tjxt/apps/message/api/internal/svc"
	"tjxt/pkg/response"
)

func ListSmsPlatformsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := smsplatform.NewListSmsPlatformsLogic(r.Context(), svcCtx)
		resp, err := l.ListSmsPlatforms()
		response.Write(w, r, resp, err)
	}
}
