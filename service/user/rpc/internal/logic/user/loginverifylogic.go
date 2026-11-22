package userlogic

import (
	"context"
	"errors"

	"tjxt/common/errorx"
	"tjxt/service/user/rpc/internal/model"
	"tjxt/service/user/rpc/internal/svc"
	"tjxt/service/user/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"golang.org/x/crypto/bcrypt"
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
//
// 失败语义统一走 pkg/errorx 业务错误，保证跨 gRPC 传输后调用方
// （api/auth 登录链路）仍能还原为 401/403 而非 500：
//   - 未提供账号标识             -> 400 请求参数错误
//   - 账号不存在 / 密码错误     -> 401 未认证
//   - 账号被禁用（status != 1） -> 403 无权限
func (l *LoginVerifyLogic) LoginVerify(in *pb.LoginVerifyRequest) (*pb.LoginVerifyResponse, error) {
	var user *model.User
	var err error

	if in.CellPhone != "" {
		user, err = l.svcCtx.UserModel.FindOneByCellPhoneType(l.ctx, in.CellPhone, int64(in.Type))
	} else if in.Username != "" {
		user, err = l.svcCtx.UserModel.FindOneByUsername(l.ctx, in.Username)
	} else {
		return nil, errorx.BadRequestf("手机号与用户名不能均为空")
	}
	if err != nil {
		if errors.Is(err, sqlc.ErrNotFound) {
			return nil, errorx.Unauthorized("账号不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询用户失败")
	}
	if user.Status != 1 {
		return nil, errorx.Forbidden("账号已被禁用")
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(in.Password)) != nil {
		return nil, errorx.Unauthorized("用户名或密码错误")
	}
	return &pb.LoginVerifyResponse{
		UserId:    user.Id,
		Username:  user.Username,
		CellPhone: user.CellPhone,
		Type:      int32(user.Type),
		Status:    int32(user.Status),
	}, nil
}
