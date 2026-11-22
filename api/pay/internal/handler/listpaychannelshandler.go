// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/pay/api/internal/logic"
	"tjxt/apps/pay/api/internal/svc"
	"tjxt/pkg/response"
)

func ListPayChannelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewListPayChannelsLogic(r.Context(), svcCtx)
		resp, err := l.ListPayChannels()
		response.Write(w, r, resp, err)
	}
}
