package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

const (
	// userColumns 账户表列列表
	userColumns = `id, username, cell_phone, password, type, status, create_time, update_time, creater, updater`

	// Cache key 模式：cache:user:id:<id>
	cacheUserKeyPrefix = "cache:user:id:"
	// Cache key 模式：cache:user:phone:<cellphone>
	cacheUserPhoneKeyPrefix = "cache:user:phone:"
	// Cache key 模式：cache:user:login:<account>:<type>
	cacheUserLoginKeyPrefix = "cache:user:login:"
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
	c     cache.Cache
	table string
}

// NewUserModel 创建账户数据访问对象，自动接入 Redis 缓存
func NewUserModel(conn sqlx.SqlConn, c cache.Cache) *UserModel {
	return &UserModel{
		conn:  conn,
		c:     c,
		table: "`user`",
	}
}

// FindById 按主键查询账户（带缓存）
// 缓存 key 模式：cache:user:id:<id>
func (m *UserModel) FindById(ctx context.Context, id int64) (*User, error) {
	cacheKey := fmt.Sprintf("%s%d", cacheUserKeyPrefix, id)
	var u User
	if err := m.c.GetCtx(ctx, cacheKey, &u); err == nil {
		return &u, nil
	}

	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, userColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err == nil && m.c != nil {
		_ = m.c.SetWithExpireCtx(ctx, cacheKey, u, 10*time.Minute)
	}
	return &u, err
}

// FindByLogin 按手机号或用户名 + 用户类型查询账户（带缓存）
// 缓存 key 模式：cache:user:login:<account>:<type>
func (m *UserModel) FindByLogin(ctx context.Context, account string, userType int64) (*User, error) {
	cacheKey := fmt.Sprintf("%s%s:%d", cacheUserLoginKeyPrefix, account, userType)
	var u User
	if err := m.c.GetCtx(ctx, cacheKey, &u); err == nil {
		return &u, nil
	}

	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE type = ? AND (cell_phone = ? OR username = ?) LIMIT 1`,
		userColumns, m.table), userType, account, account)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err == nil && m.c != nil {
		_ = m.c.SetWithExpireCtx(ctx, cacheKey, u, 10*time.Minute)
	}
	return &u, err
}

// FindByCellPhone 按手机号查询账户（带缓存）
// 缓存 key 模式：cache:user:phone:<cellphone>
func (m *UserModel) FindByCellPhone(ctx context.Context, cellPhone string) (*User, error) {
	cacheKey := fmt.Sprintf("%s%s", cacheUserPhoneKeyPrefix, cellPhone)
	var u User
	if err := m.c.GetCtx(ctx, cacheKey, &u); err == nil {
		return &u, nil
	}

	err := m.conn.QueryRowCtx(ctx, &u, fmt.Sprintf(
		`SELECT %s FROM %s WHERE cell_phone = ? LIMIT 1`, userColumns, m.table), cellPhone)
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err == nil && m.c != nil {
		_ = m.c.SetWithExpireCtx(ctx, cacheKey, u, 10*time.Minute)
	}
	return &u, err
}

// FindListByIds 批量查询账户
func (m *UserModel) FindListByIds(ctx context.Context, ids []int64) ([]*User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var users []*User
	err := m.conn.QueryRowsCtx(ctx, &users, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id IN (%s)`, userColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// Insert 新增账户（清除缓存）
func (m *UserModel) Insert(ctx context.Context, u *User) error {
	// 清除相关缓存
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheUserKeyPrefix, u.Id))
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%s", cacheUserPhoneKeyPrefix, u.CellPhone))
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%s:%d", cacheUserLoginKeyPrefix, u.Username, u.Type))
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%s:%d", cacheUserLoginKeyPrefix, u.CellPhone, u.Type))
	}

	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, username, cell_phone, password, type, status, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.table),
		u.Id, u.Username, u.CellPhone, u.Password, u.Type, u.Status)
	return err
}

// InsertTx 事务内新增账户（不触发缓存清除）
func (m *UserModel) InsertTx(ctx context.Context, tx sqlx.Session, u *User) error {
	_, err := tx.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, username, cell_phone, password, type, status, create_time, update_time) VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`,
		m.table),
		u.Id, u.Username, u.CellPhone, u.Password, u.Type, u.Status)
	return err
}

// TransactCtx 事务执行入口
func (m *UserModel) TransactCtx(ctx context.Context, fn func(ctx context.Context, tx sqlx.Session) error) error {
	return m.conn.TransactCtx(ctx, fn)
}

// UpdatePassword 修改密码（清除缓存）
func (m *UserModel) UpdatePassword(ctx context.Context, id int64, password string) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheUserKeyPrefix, id))
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET password = ?, update_time = NOW() WHERE id = ?`, m.table), password, id)
	return err
}

// UpdateStatus 修改账户状态（清除缓存）
func (m *UserModel) UpdateStatus(ctx context.Context, id, status int64) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheUserKeyPrefix, id))
	}
	_, err := m.conn.ExecCtx(ctx, fmt.Sprintf(
		`UPDATE %s SET status = ?, update_time = NOW() WHERE id = ?`, m.table), status, id)
	return err
}

// UpdateAccount 更新用户名/手机号（清除缓存）
func (m *UserModel) UpdateAccount(ctx context.Context, id int64, username, cellPhone string) error {
	if m.c != nil {
		_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%d", cacheUserKeyPrefix, id))
		if cellPhone != "" {
			_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%s", cacheUserPhoneKeyPrefix, cellPhone))
		}
		if username != "" {
			_ = m.c.DelCtx(ctx, fmt.Sprintf("%s%s:%d", cacheUserLoginKeyPrefix, username, 0))
		}
	}

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