package authlogic

import (
	"context"

	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuTreeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetMenuTreeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuTreeLogic {
	return &GetMenuTreeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetMenuTree 查询完整菜单树。
// 菜单总量有限，一次性取出后在内存中建树，避免递归查库产生 N+1。
func (l *GetMenuTreeLogic) GetMenuTree(in *pb.Empty) (*pb.MenuTreeReply, error) {
	menus, err := l.svcCtx.MenuModel.FindAll(l.ctx)
	if err != nil {
		return nil, err
	}
	return &pb.MenuTreeReply{List: buildMenuTree(menus)}, nil
}
