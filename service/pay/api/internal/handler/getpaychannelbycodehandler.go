// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/service/pay/api/internal/logic"
	"tjxt/service/pay/api/internal/svc"
	"tjxt/service/pay/api/internal/types"
	"tjxt/common/response"
)

func GetPayChannelByCodeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PayChannelCodeReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGetPayChannelByCodeLogic(r.Context(), svcCtx)
		resp, err := l.GetPayChannelByCode(&req)
		response.Write(w, r, resp, err)
	}
}
