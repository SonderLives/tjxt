package model

import (
	"context"
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
}

// AuthModel 认证相关数据访问（登录记录 + 角色名查询）
type AuthModel struct {
	conn      sqlx.SqlConn
	loginTbl  string
	roleTable string
}

// NewAuthModel 创建认证数据访问对象。
// loginTable 为登录记录表名（默认 tj_auth.login_record），roleTable 为角色表名（默认 tj_auth.role）。
func NewAuthModel(conn sqlx.SqlConn, loginTable, roleTable string) *AuthModel {
	if loginTable == "" {
		loginTable = "tj_auth.login_record"
	}
	if roleTable == "" {
		roleTable = "tj_auth.role"
	}
	return &AuthModel{conn: conn, loginTbl: loginTable, roleTable: roleTable}
}

// InsertLoginRecord 记录登录成功日志（尽力而为，失败不影响登录）。
func (m *AuthModel) InsertLoginRecord(ctx context.Context, userId int64, cellPhone string) {
	now := time.Now()
	_, _ = m.conn.ExecCtx(ctx, fmt.Sprintf(
		`INSERT INTO %s (id, user_id, cell_phone, login_time, login_date)
		 VALUES (?, ?, ?, ?, ?)`, m.loginTbl),
		idgen.NextID(), userId, cellPhone, now, now)
}

// RoleName 根据角色 id 查询角色名称。
func (m *AuthModel) RoleName(ctx context.Context, roleId int64) string {
	if roleId <= 0 {
		return ""
	}
	var name string
	err := m.conn.QueryRowCtx(ctx, &name, fmt.Sprintf(
		`SELECT name FROM %s WHERE id = ? AND deleted = 0 LIMIT 1`, m.roleTable), roleId)
	if err != nil {
		return ""
	}
	return name
}
