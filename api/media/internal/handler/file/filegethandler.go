// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package file

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/apps/media/api/internal/logic/file"
	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
	"tjxt/pkg/response"
)

func FileGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileIdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := file.NewFileGetLogic(r.Context(), svcCtx)
		resp, err := l.FileGet(&req)
		response.Write(w, r, resp, err)
	}
}
