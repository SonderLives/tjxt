// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package media

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"tjxt/service/media/api/internal/logic/media"
	"tjxt/service/media/api/internal/svc"
	"tjxt/service/media/api/internal/types"
	"tjxt/common/response"
)

func MediaSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MediaSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := media.NewMediaSaveLogic(r.Context(), svcCtx)
		resp, err := l.MediaSave(&req)
		response.Write(w, r, resp, err)
	}
}
