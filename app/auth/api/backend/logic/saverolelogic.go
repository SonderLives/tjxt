package authlogic

import (
	"context"
	"strings"

	"tjxt/app/auth/api/backend/model"
	"tjxt/app/auth/api/backend/pb"
	"tjxt/app/auth/api/backend/svc"
	"tjxt/pkg/errorx"
	"tjxt/pkg/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveRoleLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveRoleLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveRoleLogic {
	return &SaveRoleLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveRole 新增或更新角色。in.Id 为 0 表示新增，否则按主键更新。
func (l *SaveRoleLogic) SaveRole(in *pb.RoleSaveReq) (*pb.IdReply, error) {
	code := strings.TrimSpace(in.Code)
	name := strings.TrimSpace(in.Name)
	if code == "" {
		return nil, errorx.BadRequestf("角色代号不能为空")
	}
	if name == "" {
		return nil, errorx.BadRequestf("角色名称不能为空")
	}

	// code 全局唯一，更新时排除自身。
	exists, err := l.svcCtx.RoleModel.ExistsByCode(l.ctx, code, in.Id)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errorx.Conflict("角色代号已存在")
	}

	if in.Id <= 0 {
		role := &model.Role{
			Id:   idgen.NextID(),
			Code: code,
			Name: name,
			Type: int64(in.Type),
		}
		if _, err := l.svcCtx.RoleModel.Insert(l.ctx, role); err != nil {
			return nil, err
		}
		return &pb.IdReply{Id: role.Id}, nil
	}

	role, err := l.svcCtx.RoleModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NotFound("角色不存在")
		}
		return nil, err
	}
	if role.Deleted != 0 {
		return nil, errorx.NotFound("角色不存在")
	}
	// 固定角色为系统内置，禁止改动代号，避免鉴权规则失效。
	if role.Type == fixedRoleType && role.Code != code {
		return nil, errorx.Conflict("固定角色不允许修改代号")
	}

	role.Code = code
	role.Name = name
	role.Type = int64(in.Type)
	if err := l.svcCtx.RoleModel.Update(l.ctx, role); err != nil {
		return nil, err
	}
	return &pb.IdReply{Id: role.Id}, nil
}
