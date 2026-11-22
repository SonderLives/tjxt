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

type UpdateUserStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserStatusLogic {
	return &UpdateUserStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// UpdateUserStatus 更新账户状态（0-禁用 1-正常）。
func (l *UpdateUserStatusLogic) UpdateUserStatus(in *pb.UpdateStatusRequest) (*pb.EmptyResponse, error) {
	if in.Status != 0 && in.Status != 1 {
		return nil, xerr.BadRequestf("状态取值非法")
	}
	if err := l.svcCtx.UserModel.UpdateStatus(l.ctx, in.UserId, int64(in.Status), in.Operator); err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, xerr.NotFound("用户不存在")
		}
		return nil, xerr.Wrap(err, xerr.CodeInternal, "更新状态失败")
	}
	return &pb.EmptyResponse{}, nil
}
