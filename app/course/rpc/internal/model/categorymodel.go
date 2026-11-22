package model

import (
	"context"
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
		// ListAll 查询全部未删除分类，按 level、priority 升序。
		ListAll(ctx context.Context) ([]*Category, error)
		// FindByParentId 查询指定父分类的直接子分类。
		FindByParentId(ctx context.Context, parentId int64) ([]*Category, error)
		// FindByLevel 查询指定级别（1/2/3）的全部分类。
		FindByLevel(ctx context.Context, level int64) ([]*Category, error)
		// FindByCondition 按名称模糊、状态精确过滤。
		FindByCondition(ctx context.Context, name string, status int64) ([]*Category, error)
		// CountByParentId 统计直接子分类数量。
		CountByParentId(ctx context.Context, parentId int64) (int64, error)
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

func (m *customCategoryModel) ListAll(ctx context.Context) ([]*Category, error) {
	var list []*Category
	query := fmt.Sprintf("select * from %s where `deleted` = 0 order by `level` asc, `priority` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCategoryModel) FindByParentId(ctx context.Context, parentId int64) ([]*Category, error) {
	var list []*Category
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `parent_id` = ? order by `priority` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, parentId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCategoryModel) FindByLevel(ctx context.Context, level int64) ([]*Category, error) {
	var list []*Category
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `level` = ? order by `priority` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, level); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCategoryModel) FindByCondition(ctx context.Context, name string, status int64) ([]*Category, error) {
	var list []*Category
	conds := []string{"`deleted` = 0"}
	args := []any{}
	if name != "" {
		conds = append(conds, "`name` like ?")
		args = append(args, "%"+name+"%")
	}
	if status != 0 {
		conds = append(conds, "`status` = ?")
		args = append(args, status)
	}
	query := fmt.Sprintf("select * from %s where %s order by `level` asc, `priority` asc", m.table, strings.Join(conds, " and "))
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCategoryModel) CountByParentId(ctx context.Context, parentId int64) (int64, error) {
	var total int64
	query := fmt.Sprintf("select count(*) from %s where `deleted` = 0 and `parent_id` = ?", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, query, parentId); err != nil {
		return 0, err
	}
	return total, nil
}
