package handler

import (
	"net/http"

	"trade/internal/logic"
	"trade/internal/svc"
	"trade/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CartListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewCartListLogic(r.Context(), svcCtx)
		resp, err := l.CartList()
		writeResult(w, r, resp, err)
	}
}

func CartAddHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CartsAddReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCartAddLogic(r.Context(), svcCtx)
		resp, err := l.CartAdd(&req)
		writeResult(w, r, resp, err)
	}
}

func CartDeleteHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CartsDeleteReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := logic.NewCartDeleteLogic(r.Context(), svcCtx)
		resp, err := l.CartDelete(&req)
		writeResult(w, r, resp, err)
	}
}
