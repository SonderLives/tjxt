// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/apps/trade/api/internal/logic"
	"tjxt/apps/trade/api/internal/svc"
	result "tjxt/pkg/response"
)

func PayChannelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewPayChannelsLogic(r.Context(), svcCtx)
		resp, err := l.PayChannels()
		result.Write(w, r, resp, err)
	}
}
