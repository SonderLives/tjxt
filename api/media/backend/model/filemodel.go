package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FileModel = (*customFileModel)(nil)

type (
	// FileModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileModel.
	FileModel interface {
		fileModel

		// SoftDelete 逻辑删除（deleted = 1），幂等
		SoftDelete(ctx context.Context, id int64) error
		// FindOneNotDeleted 查询未删除的文件
		FindOneNotDeleted(ctx context.Context, id int64) (*File, error)
		// FindByKey 按云端 key 查询未删除的文件
		FindByKey(ctx context.Context, key string) (*File, error)
		// UpdateStatus 更新文件状态，并清理缓存
		UpdateStatus(ctx context.Context, id, status int64) error
	}

	customFileModel struct {
		*defaultFileModel
	}
)

// NewFileModel returns a model for the database table.
func NewFileModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FileModel {
	return &customFileModel{
		defaultFileModel: newFileModel(conn, c, opts...),
	}
}

// SoftDelete 逻辑删除：将 deleted 置为 1，并清理缓存
func (m *customFileModel) SoftDelete(ctx context.Context, id int64) error {
	fileIdKey := fmt.Sprintf("%s%v", cacheFileIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `deleted` = 1 where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, id)
	}, fileIdKey)
	return err
}

// FindOneNotDeleted 按 id 查询未删除的文件
func (m *customFileModel) FindOneNotDeleted(ctx context.Context, id int64) (*File, error) {
	fileIdKey := fmt.Sprintf("%s%v", cacheFileIdPrefix, id)
	var resp File
	err := m.QueryRowCtx(ctx, &resp, fileIdKey, func(ctx context.Context, conn sqlx.SqlConn, v any) error {
		query := fmt.Sprintf("select %s from %s where `id` = ? and `deleted` = 0 limit 1", fileRows, m.table)
		return conn.QueryRowCtx(ctx, v, query, id)
	})
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByKey 按云端 key 查询未删除的文件
func (m *customFileModel) FindByKey(ctx context.Context, key string) (*File, error) {
	var resp File
	query := fmt.Sprintf("select %s from %s where `key` = ? and `deleted` = 0 limit 1", fileRows, m.table)
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, key)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

// UpdateStatus 更新文件状态，并清理缓存
func (m *customFileModel) UpdateStatus(ctx context.Context, id, status int64) error {
	fileIdKey := fmt.Sprintf("%s%v", cacheFileIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set `status` = ? where `id` = ?", m.table)
		return conn.ExecCtx(ctx, query, status, id)
	}, fileIdKey)
	return err
}
