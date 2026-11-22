// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"net/http"

	"tjxt/api/course/internal/logic"
	"tjxt/api/course/internal/svc"
	"tjxt/pkg/response"
)

func CategoryListOneLevelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := logic.NewCategoryListOneLevelLogic(r.Context(), svcCtx)
		resp, err := l.CategoryListOneLevel()
		response.Write(w, r, resp, err)
	}
}
