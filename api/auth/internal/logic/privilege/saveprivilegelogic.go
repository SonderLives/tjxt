// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package privilege

import (
	"context"

	backendlogic "tjxt/api/auth/backend/logic"
	authpb "tjxt/api/auth/backend/pb"
	"tjxt/api/auth/internal/svc"
	"tjxt/api/auth/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePrivilegeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSavePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePrivilegeLogic {
	return &SavePrivilegeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SavePrivilege 新增或修改权限，返回权限 id。
func (l *SavePrivilegeLogic) SavePrivilege(req *types.PrivilegeSaveReq) (resp *types.IdVO, err error) {
	reply, err := backendlogic.NewSavePrivilegeLogic(l.ctx, l.svcCtx.Backend).SavePrivilege(&authpb.PrivilegeSaveReq{
		Id:       req.Id,
		MenuId:   req.MenuId,
		Intro:    req.Intro,
		Method:   req.Method,
		Uri:      req.Uri,
		Internal: req.Internal,
	})
	if err != nil {
		return nil, err
	}
	return &types.IdVO{Id: reply.Id}, nil
}
