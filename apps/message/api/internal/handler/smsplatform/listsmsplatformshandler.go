// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package smsplatform

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/message/api/internal/logic/smsplatform"
	"tjxt/apps/message/api/internal/svc"
)

func ListSmsPlatformsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := smsplatform.NewListSmsPlatformsLogic(r.Context(), svcCtx)
		resp, err := l.ListSmsPlatforms()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
