package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ AccountRoleModel = (*customAccountRoleModel)(nil)

type (
	// AccountRoleModel is an interface to be customized, add more methods here,
	// and implement the added methods in customAccountRoleModel.
	AccountRoleModel interface {
		accountRoleModel
		// FindRoleIdsByAccountId 查询账户已分配的角色 id 列表。
		FindRoleIdsByAccountId(ctx context.Context, accountId int64) ([]int64, error)
		// ReplaceByAccountId 全量覆盖账户的角色分配：事务内先清空再批量写入。
		ReplaceByAccountId(ctx context.Context, accountId int64, roleIds []int64) error
		// DeleteByRoleId 清空某角色的全部账户分配，用于删除角色时级联清理。
		DeleteByRoleId(ctx context.Context, roleId int64) error
		// CountByRoleId 统计仍在使用该角色的账户数，用于删除角色前的占用校验。
		CountByRoleId(ctx context.Context, roleId int64) (int64, error)
	}

	customAccountRoleModel struct {
		*defaultAccountRoleModel
	}
)

// NewAccountRoleModel returns a model for the database table.
func NewAccountRoleModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) AccountRoleModel {
	return &customAccountRoleModel{
		defaultAccountRoleModel: newAccountRoleModel(conn, c, opts...),
	}
}

// FindRoleIdsByAccountId 查询账户的角色 id 列表。
func (m *customAccountRoleModel) FindRoleIdsByAccountId(ctx context.Context, accountId int64) ([]int64, error) {
	query := fmt.Sprintf("select `role_id` from %s where `account_id` = ?", m.table)
	var ids []int64
	err := m.QueryRowsNoCacheCtx(ctx, &ids, query, accountId)
	return ids, err
}

// ReplaceByAccountId 全量覆盖账户角色分配。
func (m *customAccountRoleModel) ReplaceByAccountId(ctx context.Context, accountId int64, roleIds []int64) error {
	roleIds = dedupeIds(roleIds)

	var staleKeys []int64
	if err := m.QueryRowsNoCacheCtx(ctx, &staleKeys, fmt.Sprintf("select `id` from %s where `account_id` = ?", m.table), accountId); err != nil && err != sql.ErrNoRows {
		return err
	}

	err := m.TransactCtx(ctx, func(ctx context.Context, session sqlx.Session) error {
		if _, err := session.ExecCtx(ctx, fmt.Sprintf("delete from %s where `account_id` = ?", m.table), accountId); err != nil {
			return err
		}
		if len(roleIds) == 0 {
			return nil
		}
		args := make([]any, 0, len(roleIds)*2)
		for _, id := range roleIds {
			args = append(args, accountId, id)
		}
		insert := fmt.Sprintf("insert into %s (`account_id`, `role_id`) values %s", m.table, pairValuePlaceholders(len(roleIds)))
		_, err := session.ExecCtx(ctx, insert, args...)
		return err
	})
	if err != nil {
		return err
	}
	return m.delKeys(ctx, staleKeys)
}

// DeleteByRoleId 清空角色的账户分配。
func (m *customAccountRoleModel) DeleteByRoleId(ctx context.Context, roleId int64) error {
	var ids []int64
	if err := m.QueryRowsNoCacheCtx(ctx, &ids, fmt.Sprintf("select `id` from %s where `role_id` = ?", m.table), roleId); err != nil && err != sql.ErrNoRows {
		return err
	}
	if _, err := m.ExecNoCacheCtx(ctx, fmt.Sprintf("delete from %s where `role_id` = ?", m.table), roleId); err != nil {
		return err
	}
	return m.delKeys(ctx, ids)
}

// CountByRoleId 统计使用该角色的账户数。
func (m *customAccountRoleModel) CountByRoleId(ctx context.Context, roleId int64) (int64, error) {
	query := fmt.Sprintf("select count(*) from %s where `role_id` = ?", m.table)
	var cnt int64
	err := m.QueryRowNoCacheCtx(ctx, &cnt, query, roleId)
	return cnt, err
}

// delKeys 批量失效主键行缓存。
func (m *customAccountRoleModel) delKeys(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := m.DelCacheCtx(ctx, fmt.Sprintf("%s%v", cacheAccountRoleIdPrefix, id)); err != nil {
			return err
		}
	}
	return nil
}
