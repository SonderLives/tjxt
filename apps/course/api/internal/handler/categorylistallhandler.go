// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"tjxt/apps/course/api/internal/logic"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"
	result "tjxt/pkg/response"
)

func CategoryListAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryListAllReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCategoryListAllLogic(r.Context(), svcCtx)
		resp, err := l.CategoryListAll(&req)
		result.Write(w, r, resp, err)
	}
}
