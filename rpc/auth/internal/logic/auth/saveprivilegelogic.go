package authlogic

import (
	"context"
	"database/sql"
	"strings"

	"tjxt/apps/auth/rpc/internal/model"
	"tjxt/apps/auth/rpc/internal/svc"
	"tjxt/apps/auth/rpc/pb"
	"tjxt/pkg/utils/idgen"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type SavePrivilegeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSavePrivilegeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SavePrivilegeLogic {
	return &SavePrivilegeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SavePrivilege 新增或更新接口权限。
func (l *SavePrivilegeLogic) SavePrivilege(in *pb.PrivilegeSaveReq) (*pb.IdReply, error) {
	method := strings.ToUpper(strings.TrimSpace(in.Method))
	uri := strings.TrimSpace(in.Uri)
	if method == "" {
		return nil, xerr.BadRequestf("请求方式不能为空")
	}
	if uri == "" {
		return nil, xerr.BadRequestf("请求路径不能为空")
	}
	if in.MenuId <= 0 {
		return nil, xerr.BadRequestf("菜单 id 无效")
	}

	// 权限必须挂在有效菜单下，否则前端菜单树无法展示。
	menu, err := l.svcCtx.MenuModel.FindOne(l.ctx, in.MenuId)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("菜单不存在")
		}
		return nil, err
	}
	if menu.Deleted != 0 {
		return nil, xerr.NotFound("菜单不存在")
	}

	if in.Id <= 0 {
		privilege := &model.Privilege{
			Id:       idgen.NextID(),
			MenuId:   sql.NullInt64{Int64: in.MenuId, Valid: true},
			Intro:    in.Intro,
			Method:   method,
			Uri:      uri,
			Internal: boolToInt(in.Internal),
		}
		if _, err := l.svcCtx.PrivilegeModel.Insert(l.ctx, privilege); err != nil {
			return nil, err
		}
		return &pb.IdReply{Id: privilege.Id}, nil
	}

	privilege, err := l.svcCtx.PrivilegeModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, xerr.NotFound("权限不存在")
		}
		return nil, err
	}
	if privilege.Deleted != 0 {
		return nil, xerr.NotFound("权限不存在")
	}

	privilege.MenuId = sql.NullInt64{Int64: in.MenuId, Valid: true}
	privilege.Intro = in.Intro
	privilege.Method = method
	privilege.Uri = uri
	privilege.Internal = boolToInt(in.Internal)
	if err := l.svcCtx.PrivilegeModel.Update(l.ctx, privilege); err != nil {
		return nil, err
	}
	return &pb.IdReply{Id: privilege.Id}, nil
}
