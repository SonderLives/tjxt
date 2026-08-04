package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// User 账户（user 表）
type User struct {
	Id         int64
	Username   string
	CellPhone  string
	Password   string
	Type       int64
	Status     int64
	CreateTime time.Time
	UpdateTime time.Time
	Creater    sql.NullInt64
	Updater    sql.NullInt64
}

// UserModel 账户数据访问
type UserModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewUserModel 创建账户数据访问对象
func NewUserModel(conn sqlx.SqlConn) *UserModel {
	return &UserModel{conn: conn, table: "`user`"}
}

const userColumns = `id, username, cell_phone, password, type, status, create_time, update_time, creater, updater`

// FindById 按主键查询账户。
func (m *UserModel) FindById(ctx context.Context, id int64) (*User, error) {
	var u User
	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, userColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &u, err
}

// FindByLogin 按手机号或用户名 + 用户类型查询账户。
func (m *UserModel) FindByLogin(ctx context.Context, account string, userType int64) (*User, error) {
	var u User
	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE type = ? AND (cell_phone = ? OR username = ?) LIMIT 1`,
		userColumns, m.table), userType, account, account)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &u, err
}

// FindByCellPhone 按手机号查询账户（不限类型）。
func (m *UserModel) FindByCellPhone(ctx context.Context, cellPhone string) (*User, error) {
	var u User
	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE cell_phone = ? LIMIT 1`, userColumns, m.table), cellPhone)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &u, err
}

// Insert 新增账户。
func (m *UserModel) Insert(ctx context.Context, u *User) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, username, cell_phone, password, type, status, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.table),
		u.Id, u.Username, u.CellPhone, u.Password, u.Type, u.Status)
	return err
}

// UpdatePassword 修改密码。
func (m *UserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET password = ?, update_time = NOW() WHERE id = ?`, m.table), password, id)
	return err
}

// UpdateStatus 修改账户状态。
func (m *UserModel) UpdateStatus(ctx context.Context, id, status int64) error {
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, update_time = NOW() WHERE id = ?`, m.table), status, id)
	return err
}

// UpdateAccount 更新用户名/手机号（仅账户表）。
func (m *UserModel) UpdateAccount(ctx context.Context, id int64, username, cellPhone string) error {
	if username == "" && cellPhone == "" {
		return nil
	}
	sets := make([]string, 0, 2)
	args := make([]any, 0, 3)
	if username != "" {
		sets = append(sets, "username = ?")
		args = append(args, username)
	}
	if cellPhone != "" {
		sets = append(sets, "cell_phone = ?")
		args = append(args, cellPhone)
	}
	args = append(args, id)
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET %s, update_time = NOW() WHERE id = ?`, m.table, strings.Join(sets, ", ")), args...)
	return err
}
