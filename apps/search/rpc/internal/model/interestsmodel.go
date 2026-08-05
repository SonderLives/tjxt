package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ InterestsModel = (*customInterestsModel)(nil)

type (
	// InterestsModel is an interface to be customized, add more methods here,
	// and implement the added methods in customInterestsModel.
	InterestsModel interface {
		interestsModel
		// Upsert 保存用户兴趣：不存在则插入，存在则更新（id 为用户主键）。
		Upsert(ctx context.Context, id int64, interests string) error
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

// Upsert 一步完成插入或更新（主键冲突时更新 interests 字段），
// 并在成功后失效对应缓存。
func (m *customInterestsModel) Upsert(ctx context.Context, id int64, interests string) error {
	interestsIdKey := fmt.Sprintf("%s%v", cacheInterestsIdPrefix, id)
	_, err := m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values (?, ?) on duplicate key update `interests` = values(`interests`)", m.table, interestsRowsExpectAutoSet)
		return conn.ExecCtx(ctx, query, id, interests)
	}, interestsIdKey)
	return err
}
