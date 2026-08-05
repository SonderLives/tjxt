package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ PrivilegeModel = (*customPrivilegeModel)(nil)

type (
	// PrivilegeModel is an interface to be customized, add more methods here,
	// and implement the added methods in customPrivilegeModel.
	PrivilegeModel interface {
		privilegeModel
		// FindByMenuId 查询某菜单下全部未删除权限。
		FindByMenuId(ctx context.Context, menuId int64) ([]*Privilege, error)
		// FindByIds 按 id 批量查询未删除权限。
		FindByIds(ctx context.Context, ids []int64) ([]*Privilege, error)
		// SoftDelete 逻辑删除权限（deleted = 1）并失效缓存。
		SoftDelete(ctx context.Context, id, updater int64) error
	}

	customPrivilegeModel struct {
		*defaultPrivilegeModel
	}
)

// NewPrivilegeModel returns a model for the database table.
func NewPrivilegeModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) PrivilegeModel {
	return &customPrivilegeModel{
		defaultPrivilegeModel: newPrivilegeModel(conn, c, opts...),
	}
}

// FindByMenuId 查询菜单下的权限列表。
func (m *customPrivilegeModel) FindByMenuId(ctx context.Context, menuId int64) ([]*Privilege, error) {
	query := fmt.Sprintf("select %s from %s where `menu_id` = ? and `deleted` = 0 order by `id` asc", privilegeRows, m.table)
	var list []*Privilege
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, menuId)
	return list, err
}

// FindByIds 按 id 批量查询权限。
func (m *customPrivilegeModel) FindByIds(ctx context.Context, ids []int64) ([]*Privilege, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	holders, args := inPlaceholders(ids)
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", privilegeRows, m.table, holders)
	var list []*Privilege
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	return list, err
}

// SoftDelete 逻辑删除权限。
func (m *customPrivilegeModel) SoftDelete(ctx context.Context, id, updater int64) error {
	query := fmt.Sprintf("update %s set `deleted` = 1, `updater` = ?, `update_time` = now() where `id` = ?", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, updater, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cachePrivilegeIdPrefix, id))
}
