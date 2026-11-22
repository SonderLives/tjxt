package userlogic

import (
	"context"
	"errors"

	"tjxt/app/user/rpc/internal/model"
	"tjxt/app/user/rpc/internal/svc"
	"tjxt/app/user/rpc/pb"
	"tjxt/pkg/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserDetailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserDetailLogic {
	return &GetUserDetailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetUserDetail 获取当前登录用户详情（含资料）。
func (l *GetUserDetailLogic) GetUserDetail(in *pb.UserIdRequest) (*pb.UserDetailVO, error) {
	u, err := l.svcCtx.UserModel.FindOne(l.ctx, in.UserId)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, errorx.NotFound("用户不存在")
		}
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询用户失败")
	}
	d, err := l.svcCtx.UserDetailModel.FindOne(l.ctx, in.UserId)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, errorx.Wrap(err, errorx.CodeInternal, "查询用户失败")
	}
	return toUserDetailVO(u, d), nil
}
