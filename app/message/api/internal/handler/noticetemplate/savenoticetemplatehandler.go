// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/message/api/internal/logic/noticetemplate"
	"tjxt/service/message/api/internal/svc"
	"tjxt/service/message/api/internal/types"
	"tjxt/common/response"
)

func SaveNoticeTemplateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NoticeTemplateSaveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := noticetemplate.NewSaveNoticeTemplateLogic(r.Context(), svcCtx)
		resp, err := l.SaveNoticeTemplate(&req)
		response.Write(w, r, resp, err)
	}
}
