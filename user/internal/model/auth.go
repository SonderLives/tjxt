package model

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"common/idgen"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// LoginRecord 登录记录（auth 库 login_record 表，跨库写入）
type LoginRecord struct {
	Id         int64
	UserId     int64
	CellPhone  string
	LoginTime  time.Time
	LogoutTime time.Time
	LoginDate  time.Time
	Duration   int64
	IPv4       string
}

// AuthModel 认证相关数据访问（登录记录 + 角色名查询）
type AuthModel struct {
	conn      sqlx.SqlConn
	loginTbl  string
	roleTable string
}

// NewAuthModel 创建认证数据访问对象
func NewAuthModel(conn sqlx.SqlConn, loginTable, roleTable string) *AuthModel {
	if loginTable == "" {
		loginTable = "tj_auth.login_record"
	}
	if roleTable == "" {
		roleTable = "tj_auth.role"
	}
	return &AuthModel{conn: conn, loginTbl: loginTable, roleTable: roleTable}
}

// InsertLoginRecord 记录登录成功日志（尽力而为，失败不影响登录）
// 参数 ipv4 为客户端 IP，允许空字符串（当无法获取时由 DB 默认值兜底）
func (m *AuthModel) InsertLoginRecord(ctx context.Context, userId int64, cellPhone, ipv4 string) {
	now := time.Now()
	_, _ = m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, user_id, cell_phone, login_time, login_date, ipv4)
		 VALUES (?, ?, ?, ?, ?, ?)`, m.loginTbl),
		idgen.NextID(), userId, cellPhone, now, now, ipv4)
}

// RoleName 根据角色 id 查询角色名称，返回 (name, error)
// 返回 "" 表示角色不存在或 roleId<=0；返回 error 表示数据库查询异常
func (m *AuthModel) RoleName(ctx context.Context, roleId int64) (string, error) {
	if roleId <= 0 {
		return "", nil
	}
	var name string
	err := m.conn.QueryRowCtx(ctx, &name, fmt.Sprintf(
		`SELECT name FROM %s WHERE id = ? AND deleted = 0 LIMIT 1`, m.roleTable), roleId)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}
	return name, nil
}