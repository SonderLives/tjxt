package admin

import (
	"context"

	"tjxt/apps/user/api/internal/logic/convert"
	"tjxt/apps/user/api/internal/svc"
	"tjxt/apps/user/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddUserLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddUserLogic {
	return &AddUserLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// AddUser 管理端新增用户。
func (l *AddUserLogic) AddUser(req *types.UserDTO) (resp *types.IdVO, err error) {
	out, err := l.svcCtx.UserRpc.AddUser(l.ctx, convert.FromUserDTO(req))
	if err != nil {
		return nil, err
	}
	return convert.ToIdVO(out), nil
}
