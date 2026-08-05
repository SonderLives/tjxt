package userlogic

import (
	"context"
	"errors"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrBadCredential = errors.New("bad credential")
	ErrUserDisabled  = errors.New("user disabled")
)

type LoginVerifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginVerifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginVerifyLogic {
	return &LoginVerifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// LoginVerify 仅做身份核验（找用户 + 比对密码 + 状态检查）。
// 角色归属由 auth 服务通过 account_role 自查；此处不签发任何令牌。
func (l *LoginVerifyLogic) LoginVerify(in *pb.LoginVerifyRequest) (*pb.LoginVerifyResponse, error) {
	var user *model.User
	var err error

	if in.CellPhone != "" {
		user, err = l.svcCtx.UserModel.FindOneByCellPhoneType(l.ctx, in.CellPhone, int64(in.Type))
	} else if in.Username != "" {
		user, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	} else {
		return nil, ErrUserNotFound
	}
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		return nil, ErrBadCredential
	}
	return &pb.LoginVerifyResponse{
		UserId:    user.Id,
		Username:  user.Username,
		CellPhone: user.CellPhone,
		Type:      int32(user.Type),
		Status:    int32(user.Status),
	}, nil
}