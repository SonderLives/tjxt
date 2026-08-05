// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	result "tjxt/pkg/response"
	"tjxt/apps/media/api/internal/logic/media"
	"tjxt/apps/media/api/internal/svc"
	"tjxt/apps/media/api/internal/types"
)

func MediaDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MediaIdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewMediaDeleteLogic(r.Context(), svcCtx)
		resp, err := l.MediaDelete(&req)
		result.Write(w, r, resp, err)
	}
}
