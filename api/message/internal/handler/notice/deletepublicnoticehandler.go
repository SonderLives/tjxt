// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/message/internal/logic/notice"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"
	"tjxt/pkg/response"
)

func DeletePublicNoticeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notice.NewDeletePublicNoticeLogic(r.Context(), svcCtx)
		resp, err := l.DeletePublicNotice(&req)
		response.Write(w, r, resp, err)
	}
}
