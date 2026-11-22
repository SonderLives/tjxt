package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ QuestionBizModel = (*customQuestionBizModel)(nil)

type (
	// QuestionBizModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionBizModel.
	QuestionBizModel interface {
		questionBizModel

		// FindPageByBizId 按业务id分页查询题目关联
		FindPageByBizId(ctx context.Context, bizId int64, offset, limit int64) ([]*QuestionBiz, int64, error)
		// DeleteByQuestionId 删除某题目下的所有业务关联
		DeleteByQuestionId(ctx context.Context, questionId int64) error
	}

	customQuestionBizModel struct {
		*defaultQuestionBizModel
	}
)

// NewQuestionBizModel returns a model for the database table.
func NewQuestionBizModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) QuestionBizModel {
	return &customQuestionBizModel{
		defaultQuestionBizModel: newQuestionBizModel(conn, c, opts...),
	}
}

// FindPageByBizId 按业务id分页查询题目关联
func (m *customQuestionBizModel) FindPageByBizId(ctx context.Context, bizId int64, offset, limit int64) ([]*QuestionBiz, int64, error) {
	var total int64
	countSQL := fmt.Sprintf("select count(1) from %s where `biz_id` = ?", m.table)
	if err := m.QueryRowNoCacheCtx(ctx, &total, countSQL, bizId); err != nil {
		return nil, 0, err
	}

	listSQL := fmt.Sprintf("select %s from %s where `biz_id` = ? order by `id` asc limit ? offset ?", questionBizRows, m.table)
	var list []*QuestionBiz
	if err := m.QueryRowsNoCacheCtx(ctx, &list, listSQL, bizId, limit, offset); err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

// DeleteByQuestionId 删除某题目下的所有业务关联（逐条走缓存删除，保证缓存一致性）
func (m *customQuestionBizModel) DeleteByQuestionId(ctx context.Context, questionId int64) error {
	listSQL := fmt.Sprintf("select %s from %s where `question_id` = ?", questionBizRows, m.table)
	var list []*QuestionBiz
	if err := m.QueryRowsNoCacheCtx(ctx, &list, listSQL, questionId); err != nil {
		return err
	}
	for _, item := range list {
		if err := m.Delete(ctx, item.Id); err != nil {
			return err
		}
	}
	return nil
}
