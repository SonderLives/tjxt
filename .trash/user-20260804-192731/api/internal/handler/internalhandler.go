package handler

import (
	"context"
	"net/http"

	"tjxt/pkg/response"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

// ============ 内部接口：获取用户（缓存） ============

type InternalGetUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalGetUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalGetUserLogic {
	return &InternalGetUserLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *InternalGetUserLogic) InternalGetUser(req *types.UserIdReq) (resp *result.R, err error) {
	user, err := l.svcCtx.UserService.GetUser(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return result.OkData(user), nil
}

// ============ 内部接口：校验手机号 ============

type InternalCheckCellphoneLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalCheckCellphoneLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalCheckCellphoneLogic {
	return &InternalCheckCellphoneLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *InternalCheckCellphoneLogic) InternalCheckCellphone(req *types.UserCheckCellphoneReq) (resp *result.R, err error) {
	exists, err := l.svcCtx.UserService.CheckCellphone(l.ctx, req.Cellphone)
	if err != nil {
		return nil, err
	}
	return result.OkData(exists), nil
}

// ============ 内部接口：批量查询用户 ============

type InternalQueryBatchRequest struct {
	IDs []int64 `json:"ids"`
}

type InternalQueryBatchLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalQueryBatchLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalQueryBatchLogic {
	return &InternalQueryBatchLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *InternalQueryBatchLogic) InternalQueryBatch(req *InternalQueryBatchRequest) (resp *result.R, err error) {
	if len(req.IDs) == 0 {
		return result.OkData([]types.UserDTO{}), nil
	}
	m, err := l.svcCtx.UserService.GetUserList(l.ctx, req.IDs)
	if err != nil {
		return nil, err
	}
	list := make([]types.UserDTO, 0, len(m))
	for _, v := range m {
		list = append(list, *v)
	}
	return result.OkData(list), nil
}

// ============ 内部接口：健康检查 ============

type InternalHealthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewInternalHealthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InternalHealthLogic {
	return &InternalHealthLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *InternalHealthLogic) InternalHealth() (resp *result.R, err error) {
	return result.OkData(map[string]interface{}{"ok": true}), nil
}

// ============ Handler 实现 ============

func InternalGetUserHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserIdReq
		id, err := ParsePathID(pathvar.Vars(r)["id"])
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		req.Id = id
		l := NewInternalGetUserLogic(r.Context(), svcCtx)
		resp, err := l.InternalGetUser(&req)
		writeResult(w, r, resp, err)
	}
}

func InternalCheckCellphoneHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserCheckCellphoneReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := NewInternalCheckCellphoneLogic(r.Context(), svcCtx)
		resp, err := l.InternalCheckCellphone(&req)
		writeResult(w, r, resp, err)
	}
}

func InternalQueryBatchHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req InternalQueryBatchRequest
		if err := httpx.ParseJsonBody(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		l := NewInternalQueryBatchLogic(r.Context(), svcCtx)
		resp, err := l.InternalQueryBatch(&req)
		writeResult(w, r, resp, err)
	}
}

func InternalHealthHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := NewInternalHealthLogic(r.Context(), svcCtx)
		resp, err := l.InternalHealth()
		writeResult(w, r, resp, err)
	}
}