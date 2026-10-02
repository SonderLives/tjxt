// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/message/api/internal/logic/messagetemplate"
	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	result "tjxt/pkg/response"
)

func DeleteMessageTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := messagetemplate.NewDeleteMessageTemplateLogic(r.Context(), svcCtx)
		resp, err := l.DeleteMessageTemplate(&req)
		result.Write(w, r, resp, err)
	}
}
