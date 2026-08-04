package handler

import (
	"net/http"

	"tjxt/apps/course/api/internal/logic"
	"tjxt/apps/course/api/internal/svc"
	"tjxt/apps/course/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 保存老师信息 ============

func CourseTeacherSaveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTeacherSaveDTO
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCourseTeacherSaveLogic(r.Context(), svcCtx)
		resp, err := l.SaveTeachers(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 查询课程相关的老师信息 ============

func CourseTeacherHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CourseTeacherQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		id := pathId(r)
		l := logic.NewCourseTeacherLogic(r.Context(), svcCtx)
		resp, err := l.Teachers(id, &req)
		writeResult(w, r, resp, err)
	}
}
