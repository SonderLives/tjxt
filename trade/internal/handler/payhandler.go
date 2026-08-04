package handler

import (
	"net/http"

	"trade/internal/logic"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func PayChannelsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewPayChannelsLogic(r.Context(), svcCtx)
		resp, err := l.PayChannels()
		writeResult(w, r, resp, err)
	}
}

func PayOrderHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayApplyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewPayOrderLogic(r.Context(), svcCtx)
		resp, err := l.PayOrder(&req)
		writeResult(w, r, resp, err)
	}
}
