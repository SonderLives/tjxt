package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ QuestionDetailModel = (*customQuestionDetailModel)(nil)

type (
	// QuestionDetailModel is an interface to be customized, add more methods here,
	// and implement the added methods in customQuestionDetailModel.
	QuestionDetailModel interface {
		questionDetailModel
	}

	customQuestionDetailModel struct {
		*defaultQuestionDetailModel
	}
)

// NewQuestionDetailModel returns a model for the database table.
func NewQuestionDetailModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) QuestionDetailModel {
	return &customQuestionDetailModel{
		defaultQuestionDetailModel: newQuestionDetailModel(conn, c, opts...),
	}
}
