package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ QuestionModel = (*customQuestionModel)(nil)

type (
	// QuestionModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionModel.
	QuestionModel interface {
		questionModel

		// FindPage 按条件分页查询题目（name 模糊、type/cateId1/cateId2/difficulty 精确过滤）
		FindPage(ctx context.Context, name string, qtype, cateId1, cateId2, difficulty, offset, limit int64) ([]*Question, int64, error)
		// FindByIds 按 id 集合批量查询题目
		FindByIds(ctx context.Context, ids []int64) ([]*Question, error)
	}

	customQuestionModel struct {
		*defaultQuestionModel
	}
)

// NewQuestionModel returns a model for the database table.
func NewQuestionModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) QuestionModel {
	return &customQuestionModel{
		defaultQuestionModel: newQuestionModel(conn, c, opts...),
	}
}

// FindPage 按条件分页查询题目
func (m *customQuestionModel) FindPage(ctx context.Context, name string, qtype, cateId1, cateId2, difficulty, offset, limit int64) ([]*Question, int64, error) {
	var (
		cond []string
		args []any
	)
	if name != "" {
		cond = append(cond, "`name` like ?")
		args = append(args, "%"+name+"%")
	}
	if qtype > 0 {
		cond = append(cond, "`type` = ?")
		args = append(args, qtype)
	}
	if cateId1 > 0 {
		cond = append(cond, "`cate_id1` = ?")
		args = append(args, cateId1)
	}
	if cateId2 > 0 {
		cond = append(cond, "`cate_id2` = ?")
		args = append(args, cateId2)
	}
	if difficulty > 0 {
		cond = append(cond, "`difficulty` = ?")
		args = append(args, difficulty)
	}
	where := ""
	if len(cond) > 0 {
		where = " where " + strings.Join(cond, " and ")
	}

	var total int64
	countSQL := fmt.Sprintf("select count(1) from %s%s", m.table, where)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countSQL, args...); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("select %s from %s%s order by `id` desc limit ? offset ?", questionRows, m.table, where)
	listArgs := make([]any, 0, len(args)+2)
	listArgs = append(listArgs, args...)
	listArgs = append(listArgs, limit, offset)
	var list []*Question
	if err := m.QueryRowsNoCacheCtx(ctx, &list, listSQL, listArgs...); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// FindByIds 按 id 集合批量查询题目
func (m *customQuestionModel) FindByIds(ctx context.Context, ids []int64) ([]*Question, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	placeholders := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	query := fmt.Sprintf("select %s from %s where `id` in (%s)", questionRows, m.table, placeholders)
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	var list []*Question
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}
