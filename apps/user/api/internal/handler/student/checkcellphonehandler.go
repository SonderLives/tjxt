// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package student

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/user/api/internal/logic/student"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"
	result "tjxt/pkg/response"
)

func CheckCellPhoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckCellPhoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := student.NewCheckCellPhoneLogic(r.Context(), svcCtx)
		resp, err := l.CheckCellPhone(&req)
		result.Write(w, r, resp, err)
	}
}
