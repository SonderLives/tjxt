// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package like

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/api/remark/internal/logic/like"
	"tjxt/api/remark/internal/svc"
	"tjxt/api/remark/internal/types"
	"tjxt/pkg/response"
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
		response.Write(w, r, resp, err)
	}
}
