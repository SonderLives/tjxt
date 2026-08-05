package userlogic

import (
	"context"
	"errors"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ResetPassword 重置为默认初始口令。
func (l *ResetPasswordLogic) ResetPassword(in *pb.UserIdRequest) (*pb.EmptyResponse, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("用户不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "重置密码失败")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(defaultInitialPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "重置密码失败")
	}
	u.Password = string(hash)
	u.Updater = in.UserId
	if err := l.svcCtx.UserModel.Update(l.ctx, u); err != nil {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "重置密码失败")
	}
	return &pb.EmptyResponse{}, nil
}
