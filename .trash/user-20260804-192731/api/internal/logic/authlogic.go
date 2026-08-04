package logic

import (
	"context"

	"tjxt/pkg/response"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

// ============ 登录 ============

type AccountLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountLoginLogic {
	return &AccountLoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AccountLoginLogic) AccountLogin(req *types.LoginFormDTO) (resp *result.R, err error) {
	token, err := l.svcCtx.AuthService.Login(l.ctx, req, false)
	if err != nil {
		return nil, err
	}
	return result.OkData(token), nil
}

// ============ 管理端登录 ============

type AdminLoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminLoginLogic {
	return &AdminLoginLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AdminLoginLogic) AdminLogin(req *types.LoginFormDTO) (resp *result.R, err error) {
	token, err := l.svcCtx.AuthService.Login(l.ctx, req, true)
	if err != nil {
		return nil, err
	}
	return result.OkData(token), nil
}

// ============ 刷新 token ============

type AccountRefreshLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
	token  string
}

func NewAccountRefreshLogic(ctx context.Context, svcCtx *svc.ServiceContext, token string) *AccountRefreshLogic {
	return &AccountRefreshLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx, token: token}
}

func (l *AccountRefreshLogic) AccountRefresh() (resp *result.R, err error) {
	token, err := l.svcCtx.AuthService.Refresh(l.ctx, l.token)
	if err != nil {
		return nil, err
	}
	return result.OkData(token), nil
}

// ============ 退出登录 ============

type AccountLogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAccountLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AccountLogoutLogic {
	return &AccountLogoutLogic{Logger: logx.WithContext(ctx), ctx: ctx, svcCtx: svcCtx}
}

func (l *AccountLogoutLogic) AccountLogout() (resp *result.R, err error) {
	// 无状态 JWT：由客户端丢弃 token；如需服务端失效可接入 Redis 黑名单
	return result.Ok(), nil
}
