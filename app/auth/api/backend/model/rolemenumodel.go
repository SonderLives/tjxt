package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleMenuModel = (*customRoleMenuModel)(nil)

type (
	// RoleMenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleMenuModel.
	RoleMenuModel interface {
		roleMenuModel
		// FindMenuIdsByRoleId 查询角色已分配的菜单 id 列表。
		FindMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
		// ReplaceByRoleId 全量覆盖角色的菜单分配：事务内先清空再批量写入。
		ReplaceByRoleId(ctx context.Context, roleId int64, menuIds []int64) error
		// DeleteByRoleId 清空角色的全部菜单分配，用于删除角色时级联清理。
		DeleteByRoleId(ctx context.Context, roleId int64) error
		// DeleteByMenuId 清空某菜单的全部角色分配，用于删除菜单时级联清理。
		DeleteByMenuId(ctx context.Context, menuId int64) error
	}

	customRoleMenuModel struct {
		*defaultRoleMenuModel
	}
)

// NewRoleMenuModel returns a model for the database table.
func NewRoleMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RoleMenuModel {
	return &customRoleMenuModel{
		defaultRoleMenuModel: newRoleMenuModel(conn, c, opts...),
	}
}

// FindMenuIdsByRoleId 查询角色的菜单 id 列表。
func (m *customRoleMenuModel) FindMenuIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("select `menu_id` from %s where `role_id` = ?", m.table)
	var ids []int64
	err := m.QueryRowsNoCacheCtx(ctx, &ids, query, roleId)
	return ids, err
}

// ReplaceByRoleId 全量覆盖角色菜单分配。
func (m *customRoleMenuModel) ReplaceByRoleId(ctx context.Context, roleId int64, menuIds []int64) error {
	menuIds = dedupeIds(menuIds)

	staleKeys, err := m.rowKeysByRoleId(ctx, roleId)
	if err != nil {
		return err
	}

	err = m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, fmt.Sprintf("delete from %s where `role_id` = ?", m.table), roleId); err != nil {
			return err
		}
		if len(menuIds) == 0 {
			return nil
		}
		args := make([]any, 0, len(menuIds)*2)
		for _, id := range menuIds {
			args = append(args, roleId, id)
		}
		insert := fmt.Sprintf("insert into %s (`role_id`, `menu_id`) values %s", m.table, pairValuePlaceholders(len(menuIds)))
		_, err := session.ExecCtx(ctx, insert, args...)
		return err
	})
	if err != nil {
		return err
	}
	return m.delKeys(ctx, staleKeys)
}

// DeleteByRoleId 清空角色的菜单分配。
func (m *customRoleMenuModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	staleKeys, err := m.rowKeysByRoleId(ctx, roleId)
	if err != nil {
		return err
	}
	if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf("delete from %s where `role_id` = ?", m.table), roleId); err != nil {
		return err
	}
	return m.delKeys(ctx, staleKeys)
}

// DeleteByMenuId 清空菜单的角色分配。
func (m *customRoleMenuModel) DeleteByMenuId(ctx context.Context, menuId int64) error {
	var ids []int64
	if err := m.QueryRowsNoCacheCtx(ctx, &ids, fmt.Sprintf("select `id` from %s where `menu_id` = ?", m.table), menuId); err != nil && err != sql.ErrNoRows {
		return err
	}
	if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf("delete from %s where `menu_id` = ?", m.table), menuId); err != nil {
		return err
	}
	return m.delKeys(ctx, ids)
}

// rowKeysByRoleId 取出角色下所有关联行的主键，用于批量失效行缓存。
func (m *customRoleMenuModel) rowKeysByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	var ids []int64
	err := m.QueryRowsNoCacheCtx(ctx, &ids, fmt.Sprintf("select `id` from %s where `role_id` = ?", m.table), roleId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ids, nil
}

// delKeys 批量失效主键行缓存。
func (m *customRoleMenuModel) delKeys(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheRoleMenuIdPrefix, id)); err != nil {
			return err
		}
	}
	return nil
}
