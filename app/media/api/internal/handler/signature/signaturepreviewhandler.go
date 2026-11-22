// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/app/media/api/internal/logic/signature"
	"tjxt/app/media/api/internal/svc"
	"tjxt/app/media/api/internal/types"
	"tjxt/pkg/response"
)

func SignaturePreviewHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignatureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := signature.NewSignaturePreviewLogic(r.Context(), svcCtx)
		resp, err := l.SignaturePreview(&req)
		response.Write(w, r, resp, err)
	}
}
