// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/service/media/api/internal/logic/signature"
	"tjxt/service/media/api/internal/svc"
	"tjxt/service/media/api/internal/types"
	"tjxt/common/response"
)

func SignaturePlayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignatureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := signature.NewSignaturePlayLogic(r.Context(), svcCtx)
		resp, err := l.SignaturePlay(&req)
		response.Write(w, r, resp, err)
	}
}
