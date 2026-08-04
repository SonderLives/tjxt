package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// UserDetail 用户详情（跨库 tj_user.user_detail 表，仅读）
type UserDetail struct {
	Id           int64
	Type         int64
	Name         string
	Icon         string
	Job          string
	Intro        string
	Photo        string
	DepId        int64
}

// UserDetailModel 用户详情跨库数据访问
type UserDetailModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewUserDetailModel 创建用户详情跨库数据访问对象。
func NewUserDetailModel(conn sqlx.SqlConn, table string) *UserDetailModel {
	if table == "" {
		table = "tj_user.user_detail"
	}
	return &UserDetailModel{conn: conn, table: table}
}

const userDetailColumns = `id, type, COALESCE(name, ''), COALESCE(icon, ''), COALESCE(job, ''),
	COALESCE(intro, ''), COALESCE(photo, ''), dep_id`

// FindByIds 按用户 id 列表批量查询老师信息。
func (m *UserDetailModel) FindByIds(ctx context.Context, ids []int64) (map[int64]*UserDetail, error) {
	result := make(map[int64]*UserDetail, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []UserDetail
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id IN (%s)`, userDetailColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	for i := range rows {
		result[rows[i].Id] = &rows[i]
	}
	return result, nil
}

// FindById 按用户 id 查询老师信息。
func (m *UserDetailModel) FindById(ctx context.Context, id int64) (*UserDetail, error) {
	var d UserDetail
	err := m.conn.QueryRowCtx(ctx, &d, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id = ?`, userDetailColumns, m.table), id)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return &d, err
}
