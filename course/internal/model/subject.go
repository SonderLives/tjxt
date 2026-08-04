package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// Subject 题目表（subject 表）
type Subject struct {
	Id          int64
	SubjectName string
}

// SubjectModel 题目数据访问
type SubjectModel struct {
	conn  sqlx.SqlConn
	table string
}

// NewSubjectModel 创建题目数据访问对象
func NewSubjectModel(conn sqlx.SqlConn) *SubjectModel {
	return &SubjectModel{conn: conn, table: "subject"}
}

const subjectColumns = `id, subject_name`

// ListByIds 按 id 列表查询题目。
func (m *SubjectModel) ListByIds(ctx context.Context, ids []int64) ([]*Subject, error) {
	if len(ids) == 0 {
		return []*Subject{}, nil
	}
	ph := make([]string, 0, len(ids))
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		ph = append(ph, "?")
		args = append(args, id)
	}
	var rows []*Subject
	err := m.conn.QueryRowsCtx(ctx, &rows, fmt.Sprintf(
		`SELECT %s FROM %s WHERE id IN (%s)`, subjectColumns, m.table, strings.Join(ph, ",")), args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
