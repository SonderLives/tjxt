package model

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InterestsModel = (*customInterestsModel)(nil)

type (
	// InterestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInterestsModel.
	InterestsModel interface {
		interestsModel
	}

	customInterestsModel struct {
		*defaultInterestsModel
	}
)

// NewInterestsModel returns a model for the database table.
func NewInterestsModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) InterestsModel {
	return &customInterestsModel{
		defaultInterestsModel: newInterestsModel(conn, c, opts...),
	}
}
