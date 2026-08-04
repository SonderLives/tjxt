package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ SubjectModel = (*customSubjectModel)(nil)

type (
	// SubjectModel is an interface to be customized, add more methods here,
	// and implement the added methods in customSubjectModel.
	SubjectModel interface {
		subjectModel
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
