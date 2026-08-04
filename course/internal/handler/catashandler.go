package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 保存章节 ============

func CourseCatasSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []*types.CataSaveDTO
		if err := httpx.Parse(r, &list); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		step := pathStep(r)
		l := logic.NewCourseCatasSaveLogic(r.Context(), svcCtx)
		resp, err := l.SaveCatas(id, step, list)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取课程的章节 ============

func CourseCatasHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseCatasQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseCatasLogic(r.Context(), svcCtx)
		resp, err := l.Catas(id, &req)
		writeResult(w, r, resp, err)
	}
}

// ============ 课程视频（保存小节媒资） ============

func CourseMediaSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []*types.CourseMediaSaveDTO
		if err := httpx.Parse(r, &list); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseMediaSaveLogic(r.Context(), svcCtx)
		resp, err := l.SaveMedia(id, list)
		writeResult(w, r, resp, err)
	}
}

// ============ 保存小节或练习中的题目 ============

func CourseSubjectsSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var list []*types.CataSubjectDTO
		if err := httpx.Parse(r, &list); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseSubjectsSaveLogic(r.Context(), svcCtx)
		resp, err := l.SaveSubjects(id, list)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取小节或练习中的题目 ============

func CourseSubjectsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseSubjectsLogic(r.Context(), svcCtx)
		resp, err := l.Subjects(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 根据课程id，查询所有章节的序号 ============

func CourseCataIndexHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseCataIndexLogic(r.Context(), svcCtx)
		resp, err := l.CataIndexList(req.Id)
		writeResult(w, r, resp, err)
	}
}

// ============ 生成练习id ============

func CourseGeneratorHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCourseGeneratorLogic(r.Context(), svcCtx)
		resp, err := l.Generator()
		writeResult(w, r, resp, err)
	}
}
