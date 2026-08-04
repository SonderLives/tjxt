package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CategoryModel = (*customCategoryModel)(nil)

type (
	// CategoryModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCategoryModel.
	CategoryModel interface {
		categoryModel
		ListAll(ctx context.Context, admin bool) ([]*Category, error)
		FindById(ctx context.Context, id int64) (*Category, error)
		ListByParent(ctx context.Context, parentId int64) ([]*Category, error)
		CountByNameSameSibling(ctx context.Context, parentId int64, name string) (int64, error)
		UpdateById(ctx context.Context, data *Category, fields ...string) error
		UpdateStatusByIDs(ctx context.Context, ids []int64, status int64) error
		ListEnabled(ctx context.Context) ([]*Category, error)
		DeleteById(ctx context.Context, id int64) error
		FindByIds(ctx context.Context, ids []int64) ([]*Category, error)
	}

	customCategoryModel struct {
		*defaultCategoryModel
	}
)

// NewCategoryModel returns a model for the database table.
func NewCategoryModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CategoryModel {
	return &customCategoryModel{
		defaultCategoryModel: newCategoryModel(conn, c, opts...),
	}
}

// ListAll 查询所有分类
func (m *customCategoryModel) ListAll(ctx context.Context, admin bool) ([]*Category, error) {
	var resp []*Category
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 order by `level` asc, `priority` asc", categoryRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// FindById 根据ID查询分类
func (m *customCategoryModel) FindById(ctx context.Context, id int64) (*Category, error) {
	return m.FindOne(ctx, id)
}

// ListByParent 查询某父分类下的子分类
func (m *customCategoryModel) ListByParent(ctx context.Context, parentId int64) ([]*Category, error) {
	var resp []*Category
	query := fmt.Sprintf("select %s from %s where `parent_id` = ? and `deleted` = 0 order by `priority` asc", categoryRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, parentId)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// CountByNameSameSibling 统计同一父分类下同名分类数量
func (m *customCategoryModel) CountByNameSameSibling(ctx context.Context, parentId int64, name string) (int64, error) {
	query := fmt.Sprintf("select count(1) from %s where `parent_id` = ? and `name` = ? and `deleted` = 0", m.table)
	var count int64
	err := m.QueryRowNoCacheCtx(ctx, &count, query, parentId, name)
	return count, err
}

// UpdateById 更新指定字段
func (m *customCategoryModel) UpdateById(ctx context.Context, data *Category, fields ...string) error {
	if len(fields) == 0 {
		return m.Update(ctx, data)
	}
	categoryIdKey := fmt.Sprintf("%s%v", cacheCategoryIdPrefix, data.Id)
	setParts := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields)+1)
	for _, f := range fields {
		setParts = append(setParts, fmt.Sprintf("`%s` = ?", f))
		switch f {
		case "name":
			args = append(args, data.Name)
		case "priority":
			args = append(args, data.Priority)
		case "status":
			args = append(args, data.Status)
		}
	}
	args = append(args, data.Id)
	query := fmt.Sprintf("update %s set %s where `id` = ?", m.table, strings.Join(setParts, ", "))
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, categoryIdKey)
	return err
}

// UpdateStatusByIDs 批量更新状态
func (m *customCategoryModel) UpdateStatusByIDs(ctx context.Context, ids []int64, status int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("update %s set `status` = ? where `id` in (%s)", m.table, placeholders)
	args := make([]any, 0, len(ids)+1)
	args = append(args, status)
	for _, id := range ids {
		args = append(args, id)
	}
	categoryIdKey := fmt.Sprintf("%s%v", cacheCategoryIdPrefix, ids[0]) // 只失效第一个的缓存，实际项目中建议全失效
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, categoryIdKey)
	return err
}

// ListEnabled 查询启用的分类
func (m *customCategoryModel) ListEnabled(ctx context.Context) ([]*Category, error) {
	var resp []*Category
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `status` = 1 order by `level` asc, `priority` asc", categoryRows, m.table)
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// DeleteById 根据ID删除
func (m *customCategoryModel) DeleteById(ctx context.Context, id int64) error {
	return m.Delete(ctx, id)
}

// FindByIds 批量根据ID查询分类
func (m *customCategoryModel) FindByIds(ctx context.Context, ids []int64) ([]*Category, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	var resp []*Category
	query := fmt.Sprintf("select %s from %s where `id` in (%s) and `deleted` = 0", categoryRows, m.table, placeholders)
	if err := m.QueryRowsNoCacheCtx(ctx, &resp, query, args...); err != nil {
		return nil, err
	}
	return resp, nil
}
