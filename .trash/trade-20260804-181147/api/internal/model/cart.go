package model

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Cart 购物车条目（cart 表）
type Cart struct {
	Id         int64
	UserId     int64
	CourseId   int64
	CoverUrl   string
	CourseName string
	Price      int64
	CreateTime time.Time
	UpdateTime time.Time
}

// CartModel 购物车数据访问
type CartModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewCartModel 创建购物车数据访问对象
func NewCartModel(conn sqlx.SqlConn) *CartModel {
	return &CartModel{conn: conn, table: "cart"}
}

// Insert 添加购物车条目。
func (m *CartModel) Insert(ctx context.Context, c *Cart) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, user_id, course_id, cover_url, course_name, price, create_time, update_time)
		 VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`, m.table),
		c.Id, c.UserId, c.CourseId, c.CoverUrl, c.CourseName, c.Price)
	return err
}

// DeleteByIds 删除指定用户的购物车条目。
func (m *CartModel) DeleteByIds(ctx context.Context, userId int64, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	placeholders := ""
	args := make([]any, 0, len(ids)+1)
	args = append(args, userId)
	for i, id := range ids {
		if i > 0 {
			placeholders += ","
		}
		placeholders += "?"
		args = append(args, id)
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`DELETE FROM %s WHERE user_id = ? AND id IN (%s)`, m.table, placeholders), args...)
	return err
}

// FindByUserCourse 查询用户是否已将课程加入购物车。
func (m *CartModel) FindByUserCourse(ctx context.Context, userId, courseId int64) (*Cart, error) {
	var c Cart
	err := m.conn.QueryRowCtx(ctx, &c, fmt.Sprintf(
		`SELECT id, user_id, course_id, cover_url, course_name, price, create_time, update_time
		 FROM %s WHERE user_id = ? AND course_id = ? LIMIT 1`, m.table), userId, courseId)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListByUser 查询用户购物车列表。
func (m *CartModel) ListByUser(ctx context.Context, userId int64) ([]Cart, error) {
	var list []Cart
	err := m.conn.QueryRowsCtx(ctx, &list, fmt.Sprintf(
		`SELECT id, user_id, course_id, cover_url, course_name, price, create_time, update_time
		 FROM %s WHERE user_id = ? ORDER BY create_time DESC`, m.table), userId)
	return list, err
}
