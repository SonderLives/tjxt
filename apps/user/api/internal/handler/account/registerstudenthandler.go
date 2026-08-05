// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package account

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/user/api/internal/logic/account"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"
)

func RegisterStudentHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := account.NewRegisterStudentLogic(r.Context(), svcCtx)
		resp, err := l.RegisterStudent(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
