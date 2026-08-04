package handler

import (
	"net/http"

	"trade/internal/logic"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func RefundApplyCreateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundApplyReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyCreateLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyCreate(&req)
		writeResult(w, r, resp, err)
	}
}

func RefundApplyApproveHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ApproveReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyApproveLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyApprove(&req)
		writeResult(w, r, resp, err)
	}
}

func RefundApplyCancelHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundCancelReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyCancelLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyCancel(&req)
		writeResult(w, r, resp, err)
	}
}

func RefundApplyDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyDetailLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyDetail(&req)
		writeResult(w, r, resp, err)
	}
}

func RefundApplyNextHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewRefundApplyNextLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyNext()
		writeResult(w, r, resp, err)
	}
}

func RefundApplyPageHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundApplyPageReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyPageLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyPage(&req)
		writeResult(w, r, resp, err)
	}
}

func RefundApplyGetHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RefundIdReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewRefundApplyGetLogic(r.Context(), svcCtx)
		resp, err := l.RefundApplyGet(&req)
		writeResult(w, r, resp, err)
	}
}
