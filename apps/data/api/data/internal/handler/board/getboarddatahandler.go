// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"net/http"

	"tjxt/apps/data/api/data/internal/logic/board"
	"tjxt/apps/data/api/data/internal/svc"
	"tjxt/apps/data/api/data/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetBoardDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BoardDataReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := board.NewGetBoardDataLogic(r.Context(), svcCtx)
		resp, err := l.GetBoardData(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
