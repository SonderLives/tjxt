// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/service/trade/api/internal/logic"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/common/response"
)

func PayChannelListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewPayChannelListLogic(r.Context(), svcCtx)
		resp, err := l.PayChannelList()
		response.Write(w, r, resp, err)
	}
}
