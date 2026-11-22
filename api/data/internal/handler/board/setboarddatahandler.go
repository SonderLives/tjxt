// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package board

import (
	"net/http"

	"tjxt/api/data/internal/logic/board"
	"tjxt/api/data/internal/svc"
	"tjxt/api/data/internal/types"
	"tjxt/pkg/response"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SetBoardDataHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.BoardDataSetReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := board.NewSetBoardDataLogic(r.Context(), svcCtx)
		resp, err := l.SetBoardData(&req)
		response.Write(w, r, resp, err)
	}
}
