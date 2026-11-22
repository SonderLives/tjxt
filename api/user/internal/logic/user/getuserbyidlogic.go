package user

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	pb "tjxt/apps/user/rpc/pb"
)

type GetUserByIdLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserByIdLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserByIdLogic {
	return &GetUserByIdLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetUserById 按 id 获取用户。
func (l *GetUserByIdLogic) GetUserById(req *types.IdPathReq) (resp *types.UserDTO, err error) {
	out, err := l.svcCtx.UserRpc.GetUserById(l.ctx, &pb.UserIdRequest{UserId: req.Id})
	if err != nil {
		return nil, err
	}
	return convert.ToUserDTO(out), nil
}
