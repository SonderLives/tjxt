package handler

import (
	"net/http"

	"tjxt/apps/trade/api/internal/logic"
	"tjxt/apps/trade/api/internal/svc"
	"tjxt/apps/trade/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func OrderDetailCourseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderDetailCourseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewOrderDetailCourseLogic(r.Context(), svcCtx)
		resp, err := l.OrderDetailCourse(&req)
		writeResult(w, r, resp, err)
	}
}

func EnrollCourseHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnrollCourseReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewEnrollCourseLogic(r.Context(), svcCtx)
		resp, err := l.EnrollCourse(&req)
		writeResult(w, r, resp, err)
	}
}

func EnrollNumHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.EnrollNumReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewEnrollNumLogic(r.Context(), svcCtx)
		resp, err := l.EnrollNum(&req)
		writeResult(w, r, resp, err)
	}
}

func OrderDetailPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderDetailPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewOrderDetailPageLogic(r.Context(), svcCtx)
		resp, err := l.OrderDetailPage(&req)
		writeResult(w, r, resp, err)
	}
}

func PurchaseInfoHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.PurchaseInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewPurchaseInfoLogic(r.Context(), svcCtx)
		resp, err := l.PurchaseInfo(&req)
		writeResult(w, r, resp, err)
	}
}

func OrderDetailGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.OrderIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewOrderDetailGetLogic(r.Context(), svcCtx)
		resp, err := l.OrderDetailGet(&req)
		writeResult(w, r, resp, err)
	}
}
