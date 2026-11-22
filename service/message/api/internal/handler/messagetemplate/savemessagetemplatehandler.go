// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package messagetemplate

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/message/internal/logic/messagetemplate"
	"tjxt/api/message/internal/svc"
	"tjxt/api/message/internal/types"
	"tjxt/pkg/response"
)

func SaveMessageTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MessageTemplateSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := messagetemplate.NewSaveMessageTemplateLogic(r.Context(), svcCtx)
		resp, err := l.SaveMessageTemplate(&req)
		response.Write(w, r, resp, err)
	}
}
