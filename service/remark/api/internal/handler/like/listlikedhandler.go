// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package like

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/remark/api/internal/logic/like"
	"tjxt/service/remark/api/internal/svc"
	"tjxt/service/remark/api/internal/types"
	"tjxt/common/response"
)

func ListLikedHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LikeListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := like.NewListLikedLogic(r.Context(), svcCtx)
		resp, err := l.ListLiked(&req)
		response.Write(w, r, resp, err)
	}
}
