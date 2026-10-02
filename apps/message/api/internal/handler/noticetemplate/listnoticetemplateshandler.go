// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetemplate

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/message/api/internal/logic/noticetemplate"
	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
	result "tjxt/pkg/response"
)

func ListNoticeTemplatesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := noticetemplate.NewListNoticeTemplatesLogic(r.Context(), svcCtx)
		resp, err := l.ListNoticeTemplates(&req)
		result.Write(w, r, resp, err)
	}
}
