// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package loginrecord

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/service/auth/api/internal/logic/loginrecord"
	"tjxt/service/auth/api/internal/svc"
	"tjxt/service/auth/api/internal/types"
	"tjxt/common/response"
)

func ListLoginRecordsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginRecordListReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := loginrecord.NewListLoginRecordsLogic(r.Context(), svcCtx)
		resp, err := l.ListLoginRecords(&req)
		response.Write(w, r, resp, err)
	}
}
