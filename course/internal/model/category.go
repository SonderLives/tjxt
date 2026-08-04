package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Category 课程分类（category 表）
type Category struct {
	Id         int64
	Name       string
	ParentId   int64
	Level      int64
	Priority   int64
	Status     int64
	CreateTime time.Time
	UpdateTime time.Time
	Creater    sql.NullInt64
	Updater    sql.NullInt64
	Deleted    int64
}

// CategoryModel 课程分类数据访问
type CategoryModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCategoryModel 创建课程分类数据访问对象
func NewCategoryModel(conn sqlx.SqlConn) *CategoryModel {
	return &CategoryModel{conn: conn, table: "category"}
}

const categoryColumns = `id, name, parent_id, level, priority, status, create_time, update_time, creater, updater, deleted`

// FindById 按主键查询分类。
func (m *CategoryModel) FindById(ctx context.Context, id int64) (*Category, error) {
	var c Category
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ? AND deleted = 0`, categoryColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &c, err
}

// FindByIds 按主键列表批量查询分类。
func (m *CategoryModel) FindByIds(ctx context.Context, ids []int64) ([]*Category, error) {
	if len(ids) == 0 {
		return []*Category{}, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*Category
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND id IN (%s)`,
		categoryColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListAll 查询全部分类（按 priority 升序，id 降序）。
func (m *CategoryModel) ListAll(ctx context.Context, ascPriority bool) ([]*Category, error) {
	var rows []*Category
	order := "update_time DESC"
	if ascPriority {
		order = "priority ASC, update_time DESC"
	}
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 ORDER BY %s`, categoryColumns, m.table, order))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListByParent 按父分类查询子分类。
func (m *CategoryModel) ListByParent(ctx context.Context, parentId int64) ([]*Category, error) {
	var rows []*Category
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND parent_id = ? ORDER BY priority ASC, id DESC`,
		categoryColumns, m.table), parentId)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// ListEnabled 查询启用状态的分类（前端分类树）。
func (m *CategoryModel) ListEnabled(ctx context.Context) ([]*Category, error) {
	var rows []*Category
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE deleted = 0 AND status = 1 ORDER BY priority ASC, id DESC`,
		categoryColumns, m.table))
	if err != nil {
		return nil, err
	}
	return rows, nil
}

// Insert 新增分类。
func (m *CategoryModel) Insert(ctx context.Context, c *Category) (int64, error) {
	res, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (name, parent_id, level, priority, status, create_time, update_time, creater, updater)
		 VALUES (?, ?, ?, ?, ?, NOW(), NOW(), ?, ?)`, m.table),
		c.Name, c.ParentId, c.Level, c.Priority, c.Status, c.Creater, c.Updater)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// UpdateById 更新指定字段（name/priority/status），0 值字段跳过。
func (m *CategoryModel) UpdateById(ctx context.Context, c *Category, fields ...string) error {
	if len(fields) == 0 {
		return nil
	}
	sets := make([]string, 0, len(fields))
	args := make([]any, 0, len(fields))
	for _, f := range fields {
		switch f {
		case "name":
			sets = append(sets, "name = ?")
			args = append(args, c.Name)
		case "priority":
			sets = append(sets, "priority = ?")
			args = append(args, c.Priority)
		case "status":
			sets = append(sets, "status = ?")
			args = append(args, c.Status)
		}
	}
	if len(sets) == 0 {
		return nil
	}
	sets = append(sets, "update_time = NOW()")
	args = append(args, c.Id)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET %s WHERE id = ? AND deleted = 0`, m.table, strings.Join(sets, ", ")), args...)
	return err
}

// UpdateStatusByIDs 批量更新分类状态（禁用/启用联动）。
func (m *CategoryModel) UpdateStatusByIDs(ctx context.Context, ids []int64, status int64) error {
	if len(ids) == 0 {
		return nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids)+1)
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	args = append(args, status)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, update_time = NOW() WHERE id IN (%s) AND deleted = 0`,
		m.table, strings.Join(ph, ",")), args...)
	return err
}

// DeleteById 逻辑删除分类。
func (m *CategoryModel) DeleteById(ctx context.Context, id int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET deleted = 1, update_time = NOW() WHERE id = ? AND deleted = 0`, m.table), id)
	return err
}

// CountByNameSameSibling 统计同一父分类下同名（或与父分类同名）的分类数量。
func (m *CategoryModel) CountByNameSameSibling(ctx context.Context, parentId int64, name string) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND name = ? AND parent_id = ?`,
		m.table), name, parentId)
	if err != nil {
		return 0, err
	}
	var self int64
	err = m.conn.QueryRowCtx(ctx, &self, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND name = ? AND id = ?`,
		m.table), name, parentId)
	if err != nil {
		return 0, err
	}
	return total + self, nil
}

// CountThirdLevel 统计所有三级分类的数量（用于前端分类树过滤）。
func (m *CategoryModel) CountThirdLevel(ctx context.Context) (int64, error) {
	var total int64
	err := m.conn.QueryRowCtx(ctx, &total, fmt.Sprintf(
		`SELECT COUNT(*) FROM %s WHERE deleted = 0 AND level = 3`, m.table))
	if err != nil {
		return 0, err
	}
	return total, nil
}
