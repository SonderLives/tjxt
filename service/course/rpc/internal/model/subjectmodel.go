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

func (m *customSubjectModel) ListByIds(ctx context.Context, ids []int64) ([]*Subject, error) {
	if len(ids) == 0 {
		return []*Subject{}, nil
	}
	var list []*Subject
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `id` in (%s)", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}
