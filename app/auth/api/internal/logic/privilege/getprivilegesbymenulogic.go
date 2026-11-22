// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	backendlogic "tjxt/app/auth/api/backend/logic"
	authpb "tjxt/app/auth/api/backend/pb"
	"tjxt/app/auth/api/internal/svc"
	"tjxt/app/auth/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPrivilegesByMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPrivilegesByMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPrivilegesByMenuLogic {
	return &GetPrivilegesByMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetPrivilegesByMenu 查询某菜单下的全部权限。
func (l *GetPrivilegesByMenuLogic) GetPrivilegesByMenu(req *types.IdPathReq) (resp *types.PrivilegeListVO, err error) {
	reply, err := backendlogic.NewGetPrivilegesByMenuLogic(l.ctx, l.svcCtx.Backend).GetPrivilegesByMenu(&authpb.IdReq{Id: req.Id})
	if err != nil {
		return nil, err
	}
	list := make([]types.PrivilegeVO, 0, len(reply.List))
	for _, v := range reply.List {
		list = append(list, types.PrivilegeVO{
			Id:       v.Id,
			MenuId:   v.MenuId,
			Intro:    v.Intro,
			Method:   v.Method,
			Uri:      v.Uri,
			Internal: v.Internal,
		})
	}
	return &types.PrivilegeListVO{List: list}, nil
}
