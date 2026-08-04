package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 内部调用-获取课程信息 ============

func CourseInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseInfoQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseInfoLogic(r.Context(), svcCtx)
		resp, err := l.Info(id, &req)
		writeResult(w, r, resp, err)
	}
}

// ============ 内部调用-课程上架入索引库信息 ============

func CourseSearchInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseSearchInfoLogic(r.Context(), svcCtx)
		resp, err := l.SearchInfo(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 内部调用-通过老师id获取课程与出题数量 ============

func CourseInfoByTeacherIdsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.TeacherIdsQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseInfoByTeacherIdsLogic(r.Context(), svcCtx)
		resp, err := l.InfoByTeacherIds(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 内部调用-媒资被引用情况 ============

func CourseMediaUseInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.MediaIdsQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseMediaUseInfoLogic(r.Context(), svcCtx)
		resp, err := l.MediaUseInfo(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 内部调用-按名称查询课程id列表 ============

func CourseIdByNameHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.NameQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseIdByNameLogic(r.Context(), svcCtx)
		resp, err := l.IdByName(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 内部调用-小节信息 ============

func CourseSectionInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseSectionInfoLogic(r.Context(), svcCtx)
		resp, err := l.SectionInfo(req.Id)
		writeResult(w, r, resp, err)
	}
}
