package handler

import (
	"net/http"

	"tjxt/pkg/response"
	"tjxt/pkg/xerr"
	"tjxt/apps/user/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

// InternalOnly 内部接口中间件：校验 X-Internal-Token header
func InternalOnly(tokens []string) rest.Middleware {
	set := make(map[string]struct{}, len(tokens))
	for _, t := range tokens {
		set[t] = struct{}{}
	}
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if _, ok := set[r.Header.Get("X-Internal-Token")]; !ok {
				result.Write(w, r, nil, xerr.Unauthorized("内部接口需要 X-Internal-Token"))
				return
			}
			next(w, r)
		}
	}
}

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	// 公开接口：登录、注册、校验手机号
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/accounts/login",
				Handler: AccountLoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/accounts/admin/login",
				Handler: AdminLoginHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/students/register",
				Handler: StudentRegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/users/checkCellphone",
				Handler: UserCheckCellphoneHandler(serverCtx),
			},
		},
	)

	// 需登录接口
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/accounts/refresh",
				Handler: AccountRefreshHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/accounts/logout",
				Handler: AccountLogoutHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/students/password",
				Handler: StudentPasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/students/page",
				Handler: StudentPageHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/teachers/page",
				Handler: TeacherPageHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/staffs/page",
				Handler: StaffPageHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/users/me",
				Handler: UserMeHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/users",
				Handler: UserCreateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/users",
				Handler: UserUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/users/:id",
				Handler: UserGetHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/users/:id",
				Handler: UserAdminUpdateHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/users/:id/password/default",
				Handler: UserResetPasswordHandler(serverCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/users/:id/status/:status",
				Handler: UserStatusHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.Auth.AccessSecret),
	)

	// 内部接口（无 JWT，但需 X-Internal-Token）
	internalMiddleware := InternalOnly(serverCtx.Config.Internal.AccessTokens)
	server.Use(internalMiddleware)
	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/internal/health",
				Handler: InternalHealthHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/internal/users/:id",
				Handler: InternalGetUserHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/internal/users/checkCellphone",
				Handler: InternalCheckCellphoneHandler(serverCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/internal/users/query_batch",
				Handler: InternalQueryBatchHandler(serverCtx),
			},
		},
		rest.WithPrefix("/internal"),
	)
}