package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 保存课程基本信息 ============

func CourseSaveBaseInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseBaseInfoSaveDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseSaveBaseInfoLogic(r.Context(), svcCtx)
		resp, err := l.SaveBaseInfo(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取课程基础信息 ============

func CourseBaseInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseBaseInfoQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseBaseInfoLogic(r.Context(), svcCtx)
		resp, err := l.BaseInfo(id, &req)
		writeResult(w, r, resp, err)
	}
}

// ============ 校验课程名称 ============

func CourseCheckNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckNameQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseCheckNameLogic(r.Context(), svcCtx)
		resp, err := l.CheckName(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 删除课程 ============

func CourseDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseDeleteLogic(r.Context(), svcCtx)
		resp, err := l.Delete(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程分页查询 ============

func CoursePageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CoursePageQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCoursePageLogic(r.Context(), svcCtx)
		resp, err := l.Page(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程简要信息列表 ============

func CourseSimpleInfoListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SimpleInfoListQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseSimpleInfoListLogic(r.Context(), svcCtx)
		resp, err := l.SimpleInfoList(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程上架 ============

func CourseUpShelfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseIdDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseUpShelfLogic(r.Context(), svcCtx)
		resp, err := l.UpShelf(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程下架 ============

func CourseDownShelfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseIdDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseDownShelfLogic(r.Context(), svcCtx)
		resp, err := l.DownShelf(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程上架前校验 ============

func CourseCheckBeforeUpShelfHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseCheckBeforeUpShelfLogic(r.Context(), svcCtx)
		resp, err := l.CheckBeforeUpShelf(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 查询课程基本信息、目录、学习进度 ============

func CourseAndCatalogHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseAndCatalogLogic(r.Context(), svcCtx)
		resp, err := l.CourseAndCatalog(req.Id)
		writeResult(w, r, resp, err)
	}
}
