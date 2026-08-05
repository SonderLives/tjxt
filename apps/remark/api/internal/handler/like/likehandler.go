// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package like

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/remark/api/internal/logic/like"
	"tjxt/apps/remark/api/internal/svc"
	"tjxt/apps/remark/api/internal/types"
)

func LikeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeRecordFormReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := like.NewLikeLogic(r.Context(), svcCtx)
		resp, err := l.Like(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
