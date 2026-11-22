package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ MenuModel = (*customMenuModel)(nil)

type (
	// MenuModel is an interface to be customized, add more methods here,
	// and implement the added methods in customMenuModel.
	MenuModel interface {
		menuModel
		// FindAll 查询全部未删除菜单，按父节点与显示顺序排序，供内存中构建菜单树。
		FindAll(ctx context.Context) ([]*Menu, error)
		// FindByIds 按 id 批量查询未删除菜单。
		FindByIds(ctx context.Context, ids []int64) ([]*Menu, error)
		// CountChildren 统计某菜单下未删除的子菜单数量。
		CountChildren(ctx context.Context, parentId int64) (int64, error)
		// SyncHasChildren 依据实际子节点数量刷新 has_children 标记。
		SyncHasChildren(ctx context.Context, parentId int64) error
		// SoftDelete 逻辑删除菜单（deleted = 1）并失效缓存。
		SoftDelete(ctx context.Context, id, updater int64) error
	}

	customMenuModel struct {
		*defaultMenuModel
	}
)

// NewMenuModel returns a model for the database table.
func NewMenuModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) MenuModel {
	return &customMenuModel{
		defaultMenuModel: newMenuModel(conn, c, opts...),
	}
}

// FindAll 查询全部未删除菜单。
func (m *customMenuModel) FindAll(ctx context.Context) ([]*Menu, error) {
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `parent_id` asc, `priority` asc, `id` asc", menuRows, m.table)
	var list []*Menu
	err := m.QueryRowsNoCacheCtx(ctx, &list, query)
	return list, err
}

// FindByIds 按 id 批量查询未删除菜单。
func (m *customMenuModel) FindByIds(ctx context.Context, ids []int64) ([]*Menu, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	holders, args := inPlaceholders(ids)
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", menuRows, m.table, holders)
	var list []*Menu
	err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...)
	return list, err
}

// CountChildren 统计子菜单数量。
func (m *customMenuModel) CountChildren(ctx context.Context, parentId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `parent_id` = ? and `deleted` = 0", m.table)
	var cnt int64
	err := m.QueryRowNoCacheCtx(ctx, &cnt, query, parentId)
	return cnt, err
}

// SyncHasChildren 刷新父菜单的 has_children 标记。parentId 为 0 表示根，无需维护。
func (m *customMenuModel) SyncHasChildren(ctx context.Context, parentId int64) error {
	if parentId <= 0 {
		return nil
	}
	cnt, err := m.CountChildren(ctx, parentId)
	if err != nil {
		return err
	}
	has := 0
	if cnt > 0 {
		has = 1
	}
	query := fmt.Sprintf("update %s set `has_children` = ? where `id` = ?", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, has, parentId); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheMenuIdPrefix, parentId))
}

// SoftDelete 逻辑删除菜单。
func (m *customMenuModel) SoftDelete(ctx context.Context, id, updater int64) error {
	query := fmt.Sprintf("update %s set `deleted` = 1, `updater` = ?, `update_time` = now() where `id` = ?", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, updater, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheMenuIdPrefix, id))
}
