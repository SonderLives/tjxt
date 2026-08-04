package handler

import (
	"net/http"

	"course/internal/logic"
	"course/internal/svc"
	"course/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

// ============ 根据章节目录批量查询基础信息 ============

func CatalogueBatchQueryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdsQuery
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCatalogueBatchQueryLogic(r.Context(), svcCtx)
		resp, err := l.BatchQuery(&req)
		writeResult(w, r, resp, err)
	}
}

// ============ 获取小节信息 ============

func CatalogueSectionInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.IdPathReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCatalogueSectionInfoLogic(r.Context(), svcCtx)
		resp, err := l.SectionInfo(req.Id)
		writeResult(w, r, resp, err)
	}
}
