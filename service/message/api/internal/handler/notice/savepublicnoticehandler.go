// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package notice

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/message/api/internal/logic/notice"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"
	"tjxt/common/response"
)

func SavePublicNoticeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PublicNoticeSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := notice.NewSavePublicNoticeLogic(r.Context(), svcCtx)
		resp, err := l.SavePublicNotice(&req)
		response.Write(w, r, resp, err)
	}
}
