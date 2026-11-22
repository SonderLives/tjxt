package user

import (
	"context"

	"tjxt/api/user/internal/logic/convert"
	"tjxt/api/user/internal/svc"
	"tjxt/api/user/internal/types"
	"tjxt/pkg/auth"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/rpc/user/pb"
)

type GetCurrentUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCurrentUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCurrentUserLogic {
	return &GetCurrentUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetCurrentUser 获取当前登录用户详情。
func (l *GetCurrentUserLogic) GetCurrentUser() (resp *types.UserDetailVO, err error) {
	uid, err := auth.UserIdFromCtx(l.ctx)
	if err != nil {
		return nil, err
	}
	out, err := l.svcCtx.UserRpc.GetUserDetail(l.ctx, &pb.UserIdRequest{UserId: uid})
	if err != nil {
		return nil, err
	}
	return convert.ToUserDetailVO(out), nil
}
