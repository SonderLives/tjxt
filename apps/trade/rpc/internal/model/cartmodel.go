package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlc"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CartModel = (*customCartModel)(nil)

type (
	CartModel interface {
		cartModel
		ListByUserId(ctx context.Context, userId int64) ([]*Cart, error)
		FindByUserIdAndCourseId(ctx context.Context, userId, courseId int64) (*Cart, error)
	}
	customCartModel struct {
		*defaultCartModel
	}
)

func NewCartModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CartModel {
	return &customCartModel{
		defaultCartModel: newCartModel(conn, c, opts...),
	}
}

// ListByUserId 查询用户购物车全部条目（按加入时间倒序）
func (m *customCartModel) ListByUserId(ctx context.Context, userId int64) ([]*Cart, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? order by `create_time` desc", cartRows, m.table)
	var resp []*Cart
	err := m.QueryRowsNoCacheCtx(ctx, &resp, query, userId)
	switch err {
	case nil:
		return resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}

// FindByUserIdAndCourseId 查询用户是否已将该课程加入购物车
func (m *customCartModel) FindByUserIdAndCourseId(ctx context.Context, userId, courseId int64) (*Cart, error) {
	query := fmt.Sprintf("select %s from %s where `user_id` = ? and `course_id` = ? limit 1", cartRows, m.table)
	var resp Cart
	err := m.QueryRowNoCacheCtx(ctx, &resp, query, userId, courseId)
	switch err {
	case nil:
		return &resp, nil
	case sqlc.ErrNotFound:
		return nil, ErrNotFound
	default:
		return nil, err
	}
}
