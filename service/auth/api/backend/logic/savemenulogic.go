package authlogic

import (
	"context"
	"strings"

	"tjxt/api/auth/backend/model"
	"tjxt/api/auth/backend/pb"
	"tjxt/api/auth/backend/svc"
	"tjxt/pkg/errorx"
	"tjxt/pkg/idgen"

	"github.com/zeromicro/go-zero/core/logx"
)

// defaultMenuPriority 与建表默认值保持一致，未指定顺序时排在末尾。
const defaultMenuPriority int64 = 127

type SaveMenuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveMenuLogic {
	return &SaveMenuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SaveMenu 新增或更新菜单，并同步父节点的 has_children 标记。
func (l *SaveMenuLogic) SaveMenu(in *pb.MenuSaveReq) (*pb.IdReply, error) {
	label := strings.TrimSpace(in.Label)
	if label == "" {
		return nil, errorx.BadRequestf("菜单名称不能为空")
	}
	if in.ParentId < 0 {
		return nil, errorx.BadRequestf("父菜单 id 无效")
	}
	if in.Id > 0 && in.Id == in.ParentId {
		return nil, errorx.BadRequestf("父菜单不能是自身")
	}

	// 指定父节点时校验其存在，避免产生游离分支。
	if in.ParentId > 0 {
		parent, err := l.svcCtx.MenuModel.FindOne(l.ctx, in.ParentId)
		if err != nil {
			if err == model.ErrNotFound {
				return nil, errorx.NotFound("父菜单不存在")
			}
			return nil, err
		}
		if parent.Deleted != 0 {
			return nil, errorx.NotFound("父菜单不存在")
		}
	}

	priority := int64(in.Priority)
	if priority <= 0 {
		priority = defaultMenuPriority
	}

	if in.Id <= 0 {
		menu := &model.Menu{
			Id:       idgen.NextID(),
			ParentId: in.ParentId,
			Label:    label,
			Path:     in.Path,
			Icon:     in.Icon,
			Priority: priority,
		}
		if _, err := l.svcCtx.MenuModel.Insert(l.ctx, menu); err != nil {
			return nil, err
		}
		if err := l.svcCtx.MenuModel.SyncHasChildren(l.ctx, in.ParentId); err != nil {
			return nil, err
		}
		return &pb.IdReply{Id: menu.Id}, nil
	}

	menu, err := l.svcCtx.MenuModel.FindOne(l.ctx, in.Id)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, errorx.NotFound("菜单不存在")
		}
		return nil, err
	}
	if menu.Deleted != 0 {
		return nil, errorx.NotFound("菜单不存在")
	}

	oldParent := menu.ParentId
	menu.ParentId = in.ParentId
	menu.Label = label
	menu.Path = in.Path
	menu.Icon = in.Icon
	menu.Priority = priority
	if err := l.svcCtx.MenuModel.Update(l.ctx, menu); err != nil {
		return nil, err
	}

	// 父节点变更时，新旧父节点的 has_children 都要重算。
	if err := l.svcCtx.MenuModel.SyncHasChildren(l.ctx, in.ParentId); err != nil {
		return nil, err
	}
	if oldParent != in.ParentId {
		if err := l.svcCtx.MenuModel.SyncHasChildren(l.ctx, oldParent); err != nil {
			return nil, err
		}
	}
	return &pb.IdReply{Id: menu.Id}, nil
}
