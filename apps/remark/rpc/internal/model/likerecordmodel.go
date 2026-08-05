package model

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ LikeRecordModel = (*customLikeRecordModel)(nil)

type (
	// LikeRecordModel is an interface to be customized, add more methods here,
	// and implement the added methods in customLikeRecordModel.
	LikeRecordModel interface {
		likeRecordModel
		// FindLikedBizIds 查询某用户对一组业务 id 中已点赞(liked=1)的 biz_id 集合。
		FindLikedBizIds(ctx context.Context, userId int64, bizType string, bizIds []int64) ([]int64, error)
	}

	customLikeRecordModel struct {
		*defaultLikeRecordModel
	}
)

// NewLikeRecordModel returns a model for the database table.
func NewLikeRecordModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) LikeRecordModel {
	return &customLikeRecordModel{
		defaultLikeRecordModel: newLikeRecordModel(conn, c, opts...),
	}
}

// FindLikedBizIds 不带缓存的批量点赞状态查询；命中行数通常远小于入参 bizIds,
// 且业务对实时性要求高，走 NoCache 查询避免缓存击穿与脏读。
func (m *customLikeRecordModel) FindLikedBizIds(ctx context.Context, userId int64, bizType string, bizIds []int64) ([]int64, error) {
	if len(bizIds) == 0 {
		return nil, nil
	}
	placeholders := make([]string, len(bizIds))
	args := make([]any, 0, len(bizIds)+2)
	args = append(args, userId, bizType)
	for i, id := range bizIds {
		placeholders[i] = "?"
		args = append(args, id)
	}
	query := fmt.Sprintf(
		"select `biz_id` from %s where `user_id`=? and `biz_type`=? and `liked`=1 and `biz_id` in (%s)",
		m.table, strings.Join(placeholders, ","))
	var rows []int64
	if err := m.CachedConn.QueryRowsNoCacheCtx(ctx, &rows, query, args...); err != nil {
		return nil, err
	}
	return rows, nil
}
