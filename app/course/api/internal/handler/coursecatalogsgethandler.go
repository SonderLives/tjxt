// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/app/course/api/internal/logic"
	"tjxt/app/course/api/internal/svc"
	"tjxt/app/course/api/internal/types"
	"tjxt/pkg/response"
)

func CourseCatalogsGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCourseCatalogsGetLogic(r.Context(), svcCtx)
		resp, err := l.CourseCatalogsGet(&req)
		response.Write(w, r, resp, err)
	}
}
