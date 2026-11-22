package authlogic

import (
	"context"

	"tjxt/service/auth/api/backend/pb"
	"tjxt/service/auth/api/backend/svc"
	"tjxt/common/errorx"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAccountRolesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetAccountRolesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAccountRolesLogic {
	return &GetAccountRolesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetAccountRoles 查询账户已分配的角色 id 列表。
func (l *GetAccountRolesLogic) GetAccountRoles(in *pb.IdReq) (*pb.IdListReply, error) {
	if in.Id <= 0 {
		return nil, errorx.BadRequestf("账户 id 无效")
	}

	ids, err := l.svcCtx.AccountRoleModel.FindRoleIdsByAccountId(l.ctx, in.Id)
	if err != nil {
		return nil, err
	}
	if ids == nil {
		ids = []int64{}
	}
	return &pb.IdListReply{Ids: ids}, nil
}
