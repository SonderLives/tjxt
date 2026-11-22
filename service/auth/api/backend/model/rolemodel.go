package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ RoleModel = (*customRoleModel)(nil)

type (
	// RoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customRoleModel.
	RoleModel interface {
		roleModel
		// FindPage 分页查询未逻辑删除的角色，按创建时间倒序。
		FindPage(ctx context.Context, offset, limit int64) ([]*Role, int64, error)
		// ExistsByCode 校验角色 code 唯一性，excludeId > 0 时排除自身（用于更新场景）。
		ExistsByCode(ctx context.Context, code string, excludeId int64) (bool, error)
		// SoftDelete 逻辑删除角色（deleted = 1）并失效缓存。
		SoftDelete(ctx context.Context, id, updater int64) error
	}

	customRoleModel struct {
		*defaultRoleModel
	}
)

// NewRoleModel returns a model for the database table.
func NewRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) RoleModel {
	return &customRoleModel{
		defaultRoleModel: newRoleModel(conn, c, opts...),
	}
}

// FindPage 分页查询角色列表。
func (m *customRoleModel) FindPage(ctx context.Context, offset, limit int64) ([]*Role, int64, error) {
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `create_time` desc limit ?, ?", roleRows, m.table)
	var list []*Role
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, offset, limit); err != nil {
		return nil, 0, err
	}

	var total int64
	countQuery := fmt.Sprintf("select count(*) from %s where `deleted` = 0", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countQuery); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// ExistsByCode 判断角色 code 是否已被占用。
func (m *customRoleModel) ExistsByCode(ctx context.Context, code string, excludeId int64) (bool, error) {
	query := fmt.Sprintf("select count(*) from %s where `code` = ? and `deleted` = 0 and `id` <> ?", m.table)
	var cnt int64
	if err := m.QueryRowNoCacheCtx(ctx, &cnt, query, code, excludeId); err != nil {
		return false, err
	}
	return cnt > 0, nil
}

// SoftDelete 逻辑删除角色。
func (m *customRoleModel) SoftDelete(ctx context.Context, id, updater int64) error {
	key := fmt.Sprintf("%s%v", cacheRoleIdPrefix, id)
	query := fmt.Sprintf("update %s set `deleted` = 1, `updater` = ?, `update_time` = now() where `id` = ?", m.table)
	if _, err := m.ExecNoCacheCtx(ctx, query, updater, id); err != nil {
		return err
	}
	return m.DelCacheCtx(ctx, key)
}
