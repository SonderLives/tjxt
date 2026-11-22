// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package signature

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/api/media/internal/logic/signature"
	"tjxt/api/media/internal/svc"
	"tjxt/api/media/internal/types"
	"tjxt/pkg/response"
)

func SignatureUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SignatureReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := signature.NewSignatureUploadLogic(r.Context(), svcCtx)
		resp, err := l.SignatureUpload(&req)
		response.Write(w, r, resp, err)
	}
}
