// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package noticetask

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/message/api/internal/logic/noticetask"
	"tjxt/apps/message/api/internal/svc"
	"tjxt/apps/message/api/internal/types"
)

func ListNoticeTasksHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := noticetask.NewListNoticeTasksLogic(r.Context(), svcCtx)
		resp, err := l.ListNoticeTasks(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
