package service

import (
	"context"
	"database/sql"
	"strings"

	"common/auth"
	"common/xerr"
	"user/internal/model"
	"user/internal/types"

	"golang.org/x/crypto/bcrypt"
)

// AuthService 认证业务接口
type AuthService interface {
	// Login 登录并签发 JWT；isAdmin 为 true 时走管理端登录（仅员工/管理员）。
	Login(ctx context.Context, req *types.LoginFormDTO, isAdmin bool) (string, error)
	// Refresh 校验现有 token 并重新签发。
	Refresh(ctx context.Context, tokenString string) (string, error)
}

type authService struct {
	userModel *model.UserModel
	authModel *model.AuthModel
	accessSec string
	accessExp int64
}

// NewAuthService 创建认证业务服务。
func NewAuthService(userModel *model.UserModel, authModel *model.AuthModel, accessSecret string, accessExpire int64) AuthService {
	return &authService{
		userModel: userModel,
		authModel: authModel,
		accessSec: accessSecret,
		accessExp: accessExpire,
	}
}

// Login 登录。
func (s *authService) Login(ctx context.Context, req *types.LoginFormDTO, isAdmin bool) (string, error) {
	if req.Type <= 0 {
		return "", xerr.BadRequestf("用户类型不能为空")
	}
	if !isValidUserType(req.Type) {
		return "", xerr.BadRequestf("非法的用户类型")
	}
	account := strings.TrimSpace(req.CellPhone)
	if account == "" {
		account = strings.TrimSpace(req.Username)
	}
	if account == "" {
		return "", xerr.BadRequestf("手机号或用户名不能为空")
	}
	if req.Password == "" {
		return "", xerr.BadRequestf("密码不能为空")
	}

	u, err := s.userModel.FindByLogin(ctx, account, req.Type)
	if err == sql.ErrNoRows {
		return "", xerr.BadRequestf("账号或密码错误")
	}
	if err != nil {
		return "", xerr.Wrap(err, xerr.CodeInternal, "查询账户失败")
	}
	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)) != nil {
		return "", xerr.BadRequestf("账号或密码错误")
	}
	if u.Status != model.UserStatusNormal {
		return "", xerr.Forbidden("账号已被禁用")
	}
	// 管理端登录仅允许员工（管理员为 type=1 的账户）
	if isAdmin && req.Type != model.UserTypeStaff {
		return "", xerr.Forbidden("非管理端账户")
	}

	role := roleFromType(u.Type)
	if isAdmin {
		role = auth.RoleAdmin
	}
	token, err := auth.Sign(u.Id, role, s.accessSec, s.accessExp)
	if err != nil {
		return "", err
	}
	s.authModel.InsertLoginRecord(ctx, u.Id, u.CellPhone)
	return token, nil
}

// Refresh 刷新 token。
func (s *authService) Refresh(ctx context.Context, tokenString string) (string, error) {
	userID, role, err := auth.Parse(tokenString, s.accessSec)
	if err != nil {
		return "", err
	}
	token, err := auth.Sign(userID, role, s.accessSec, s.accessExp)
	if err != nil {
		return "", err
	}
	return token, nil
}

// isValidUserType 校验用户类型。
func isValidUserType(t int64) bool {
	switch t {
	case model.UserTypeStaff, model.UserTypeStudent, model.UserTypeTeacher:
		return true
	}
	return false
}

// roleFromType 由用户类型推导角色。
func roleFromType(t int64) string {
	switch t {
	case model.UserTypeStaff:
		return auth.RoleStaff
	case model.UserTypeTeacher:
		return auth.RoleTeacher
	case model.UserTypeStudent:
		return auth.RoleStudent
	}
	return auth.RoleUser
}
