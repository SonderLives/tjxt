// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/trade/api/internal/logic"
	"tjxt/service/trade/api/internal/svc"
	"tjxt/service/trade/api/internal/types"
	"tjxt/common/response"
)

func PayChannelDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayChannelIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewPayChannelDeleteLogic(r.Context(), svcCtx)
		resp, err := l.PayChannelDelete(&req)
		response.Write(w, r, resp, err)
	}
}
