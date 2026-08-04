package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubjectModel = (*customSubjectModel)(nil)

type (
	// SubjectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubjectModel.
	SubjectModel interface {
		subjectModel
		ListByIds(ctx context.Context, ids []int64) ([]*Subject, error)
	}

	customSubjectModel struct {
		*defaultSubjectModel
	}
)

// NewSubjectModel returns a model for the database table.
func NewSubjectModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) SubjectModel {
	return &customSubjectModel{
		defaultSubjectModel: newSubjectModel(conn, c, opts...),
	}
}

// ListByIds 批量查询题目
func (m *customSubjectModel) ListByIds(ctx context.Context, ids []int64) ([]*Subject, error) {
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
	query := fmt.Sprintf("select %s from %s where `deleted` = 0 and `id` in (%s)", subjectRows, m.table, strings.Join(ph, ","))
	err := m.QueryRowsNoCacheCtx(ctx, &rows, query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
