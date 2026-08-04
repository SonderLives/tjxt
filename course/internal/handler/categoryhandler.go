package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 查询课程分类信息 ============

func CategoryListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryListQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryListLogic(r.Context(), svcCtx)
		resp, err := l.CategoryList(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 新增课程分类 ============

func CategoryAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryAddDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryAddLogic(r.Context(), svcCtx)
		resp, err := l.CategoryAdd(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取课程分类信息 ============

func CategoryGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryGetLogic(r.Context(), svcCtx)
		resp, err := l.CategoryGet(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 删除分类信息 ============

func CategoryDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryDeleteLogic(r.Context(), svcCtx)
		resp, err := l.CategoryDelete(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程分类停用或启用 ============

func CategoryDisableOrEnableHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryDisableOrEnableDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryDisableOrEnableLogic(r.Context(), svcCtx)
		resp, err := l.CategoryDisableOrEnable(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 更新课程分类 ============

func CategoryUpdateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CategoryUpdateDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryUpdateLogic(r.Context(), svcCtx)
		resp, err := l.CategoryUpdate(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取所有课程分类（分类树） ============

func CategoryAllHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.AdminQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCategoryAllLogic(r.Context(), svcCtx)
		resp, err := l.CategoryAll(req.Admin)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取所有课程分类（不分层） ============

func CategoryAllOfOneLevelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCategoryAllOfOneLevelLogic(r.Context(), svcCtx)
		resp, err := l.CategoryAllOfOneLevel()
		writeResult(w, r, resp, err)
	}
}
