package userlogic

import (
	"context"
	"errors"

	"tjxt/apps/user/rpc/internal/model"
	"tjxt/apps/user/rpc/internal/svc"
	"tjxt/apps/user/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserByIdLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserById 按 id 获取用户（含资料，资料缺失时退化为仅基础信息）。
func (l *GetUserByIdLogic) GetUserById(in *pb.UserIdRequest) (*pb.UserDTO, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("用户不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户失败")
	}
	d, err := l.svcCtx.UserDetailModel.FindOne(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, xerr.Wrap(err, xerr.CodeInternal, "查询用户失败")
	}
	return toUserDTO(u, d), nil
}
