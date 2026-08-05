package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RolePrivilegeModel = (*customRolePrivilegeModel)(nil)

type (
	// RolePrivilegeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRolePrivilegeModel.
	RolePrivilegeModel interface {
		rolePrivilegeModel
		// FindPrivilegeIdsByRoleId 查询角色已分配的权限 id 列表。
		FindPrivilegeIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error)
		// ReplaceByRoleId 全量覆盖角色的权限分配：事务内先清空再批量写入。
		ReplaceByRoleId(ctx context.Context, roleId int64, privilegeIds []int64) error
		// DeleteByRoleId 清空角色的全部权限分配，用于删除角色时级联清理。
		DeleteByRoleId(ctx context.Context, roleId int64) error
		// DeleteByPrivilegeId 清空某权限的全部角色分配，用于删除权限时级联清理。
		DeleteByPrivilegeId(ctx context.Context, privilegeId int64) error
	}

	customRolePrivilegeModel struct {
		*defaultRolePrivilegeModel
	}
)

// NewRolePrivilegeModel returns a model for the database table.
func NewRolePrivilegeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RolePrivilegeModel {
	return &customRolePrivilegeModel{
		defaultRolePrivilegeModel: newRolePrivilegeModel(conn, c, opts...),
	}
}

// FindPrivilegeIdsByRoleId 查询角色的权限 id 列表。
func (m *customRolePrivilegeModel) FindPrivilegeIdsByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	query := fmt.Sprintf("select `privilege_id` from %s where `role_id` = ?", m.table)
	var ids []int64
	err := m.QueryRowsNoCacheCtx(ctx, &ids, query, roleId)
	return ids, err
}

// ReplaceByRoleId 全量覆盖角色权限分配。
func (m *customRolePrivilegeModel) ReplaceByRoleId(ctx context.Context, roleId int64, privilegeIds []int64) error {
	privilegeIds = dedupeIds(privilegeIds)

	staleKeys, err := m.rowKeysByRoleId(ctx, roleId)
	if err != nil {
		return err
	}

	err = m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, fmt.Sprintf("delete from %s where `role_id` = ?", m.table), roleId); err != nil {
			return err
		}
		if len(privilegeIds) == 0 {
			return nil
		}
		args := make([]any, 0, len(privilegeIds)*2)
		for _, id := range privilegeIds {
			args = append(args, roleId, id)
		}
		insert := fmt.Sprintf("insert into %s (`role_id`, `privilege_id`) values %s", m.table, pairValuePlaceholders(len(privilegeIds)))
		_, err := session.ExecCtx(ctx, insert, args...)
		return err
	})
	if err != nil {
		return err
	}
	return m.delKeys(ctx, staleKeys)
}

// DeleteByRoleId 清空角色的权限分配。
func (m *customRolePrivilegeModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	staleKeys, err := m.rowKeysByRoleId(ctx, roleId)
	if err != nil {
		return err
	}
	if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf("delete from %s where `role_id` = ?", m.table), roleId); err != nil {
		return err
	}
	return m.delKeys(ctx, staleKeys)
}

// DeleteByPrivilegeId 清空权限的角色分配。
func (m *customRolePrivilegeModel) DeleteByPrivilegeId(ctx context.Context, privilegeId int64) error {
	var ids []int64
	if err := m.QueryRowsNoCacheCtx(ctx, &ids, fmt.Sprintf("select `id` from %s where `privilege_id` = ?", m.table), privilegeId); err != nil && err != sql.ErrNoRows {
		return err
	}
	if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf("delete from %s where `privilege_id` = ?", m.table), privilegeId); err != nil {
		return err
	}
	return m.delKeys(ctx, ids)
}

// rowKeysByRoleId 取出角色下所有关联行的主键，用于批量失效行缓存。
func (m *customRolePrivilegeModel) rowKeysByRoleId(ctx context.Context, roleId int64) ([]int64, error) {
	var ids []int64
	err := m.QueryRowsNoCacheCtx(ctx, &ids, fmt.Sprintf("select `id` from %s where `role_id` = ?", m.table), roleId)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	return ids, nil
}

// delKeys 批量失效主键行缓存。
func (m *customRolePrivilegeModel) delKeys(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheRolePrivilegeIdPrefix, id)); err != nil {
			return err
		}
	}
	return nil
}
