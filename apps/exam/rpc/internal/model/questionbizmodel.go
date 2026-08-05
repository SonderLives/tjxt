package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ QuestionBizModel = (*customQuestionBizModel)(nil)

type (
	// QuestionBizModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionBizModel.
	QuestionBizModel interface {
		questionBizModel
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
