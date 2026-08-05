package account

import (
	"context"

	authclient "tjxt/apps/auth/rpc/client/auth"
	userclient "tjxt/apps/user/rpc/client/user"
	authutil "tjxt/pkg/auth"

	"tjxt/apps/auth/api/internal/svc"
	"tjxt/apps/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// Login 用户登录：调 user 服务核验身份，再调 auth 服务签发令牌。
func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginVO, err error) {
	verify, err := l.svcCtx.UserRpc.LoginVerify(l.ctx, &userclient.LoginVerifyRequest{
		CellPhone: req.CellPhone,
		Username:  req.Username,
		Password:  req.Password,
		Type:      req.Type,
	})
	if err != nil {
		return nil, err
	}

	// 通过 auth 库查账户角色；学员/老师若无 employee 账号，role 取 USER/STUDENT/TEACHER
	roleCode := roleByType(verify.Type)
	signed, err := l.svcCtx.AuthRpc.SignToken(l.ctx, &authclient.SignTokenReq{
		UserId:   verify.UserId,
		RoleCode: roleCode,
	})
	if err != nil {
		return nil, err
	}
	return &types.LoginVO{
		AccessToken:  signed.AccessToken,
		RefreshToken: signed.RefreshToken,
		ExpiresAt:    signed.ExpiresAt,
	}, nil
}

// roleByType 按 user.Type 给一个初始角色；真实项目里应由 account_role 表驱动。
func roleByType(userType int32) string {
	switch userType {
	case 1:
		return authutil.RoleStaff
	case 2:
		return authutil.RoleStudent
	case 3:
		return authutil.RoleTeacher
	default:
		return authutil.RoleUser
	}
}