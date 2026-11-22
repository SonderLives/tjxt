// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/api/media/internal/logic/media"
	"tjxt/api/media/internal/svc"
	"tjxt/api/media/internal/types"
	"tjxt/pkg/response"
)

func MediaListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MediaListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewMediaListLogic(r.Context(), svcCtx)
		resp, err := l.MediaList(&req)
		response.Write(w, r, resp, err)
	}
}
