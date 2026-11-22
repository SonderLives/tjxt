package model

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCatalogueModel = (*customCourseCatalogueModel)(nil)

// MediaUseCount 媒资在课程目录中的引用计数（用于 CourseMediaUseInfo）。
type MediaUseCount struct {
	MediaId  int64 `db:"media_id"`
	QuoteNum int64 `db:"quote_num"`
}

type (
	// CourseCatalogueModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCatalogueModel.
	CourseCatalogueModel interface {
		courseCatalogueModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogue, error)
		ListByIdIn(ctx context.Context, ids []int64) ([]*CourseCatalogue, error)
		CountByMediaIds(ctx context.Context, mediaIds []int64) ([]*MediaUseCount, error)
		DeleteByIdsTx(ctx context.Context, session sqlx.Session, ids []int64) error
		BatchInsertTx(ctx context.Context, session sqlx.Session, list []*CourseCatalogue) error
	}
	customCourseCatalogueModel struct {
		*defaultCourseCatalogueModel
	}
)

// NewCourseCatalogueModel returns a model for the database table.
func NewCourseCatalogueModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCatalogueModel {
	return &customCourseCatalogueModel{
		defaultCourseCatalogueModel: newCourseCatalogueModel(conn, c, opts...),
	}
}

func (m *customCourseCatalogueModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCatalogue, error) {
	var list []*CourseCatalogue
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `course_id` = ? order by `c_index` asc", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCatalogueModel) ListByIdIn(ctx context.Context, ids []int64) ([]*CourseCatalogue, error) {
	if len(ids) == 0 {
		return []*CourseCatalogue{}, nil
	}
	var list []*CourseCatalogue
	ph := strings.TrimSuffix(strings.Repeat("?,", len(ids)), ",")
	args := make([]any, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
	}
	query := fmt.Sprintf("select * from %s where `deleted` = 0 and `id` in (%s) order by `c_index` asc", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCatalogueModel) CountByMediaIds(ctx context.Context, mediaIds []int64) ([]*MediaUseCount, error) {
	if len(mediaIds) == 0 {
		return []*MediaUseCount{}, nil
	}
	var list []*MediaUseCount
	ph := strings.TrimSuffix(strings.Repeat("?,", len(mediaIds)), ",")
	args := make([]any, 0, len(mediaIds))
	for _, id := range mediaIds {
		args = append(args, id)
	}
	query := fmt.Sprintf("select `media_id`, count(*) as `quote_num` from %s where `deleted` = 0 and `media_id` in (%s) group by `media_id`", m.table, ph)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
		return nil, err
	}
	return list, nil
}

// DeleteByIdsTx 事务内按 id 批量物理删除目录行（一条 DELETE ... WHERE id IN），并失效对应行缓存，
// 供跨模型事务（如课程上架发布时清理旧目录）复用；ids 为空时为空操作。
func (m *customCourseCatalogueModel) DeleteByIdsTx(ctx context.Context, session sqlx.Session, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	query := fmt.Sprintf("delete from %s where `id` in (%s)", m.table, strings.TrimSuffix(strings.Repeat("?,", len(ids)), ","))
	args := make([]any, 0, len(ids))
	keys := make([]string, 0, len(ids))
	for _, id := range ids {
		args = append(args, id)
		keys = append(keys, fmt.Sprintf("%s%v", cacheCourseCatalogueIdPrefix, id))
	}
	cc := m.WithSession(session)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)
	return err
}

// BatchInsertTx 事务内批量插入目录行（单条 INSERT 多 VALUES，列与生成模型 Insert 一致），
// 并逐行失效行缓存，供跨模型事务（如课程上架发布）复用；list 为空时为空操作。
func (m *customCourseCatalogueModel) BatchInsertTx(ctx context.Context, session sqlx.Session, list []*CourseCatalogue) error {
	if len(list) == 0 {
		return nil
	}
	cols := insertColumns(courseCatalogueRowsExpectAutoSet)
	values := make([]string, 0, len(list))
	args := make([]any, 0, len(list)*cols)
	keys := make([]string, 0, len(list))
	for _, d := range list {
		values = append(values, rowPlaceholders(cols))
		args = append(args, d.Id, d.Name, d.Trailer, d.CourseId, d.Type, d.ParentCatalogueId,
			d.MediaId, d.VideoId, d.VideoName, d.LivingStartTime, d.LivingEndTime, d.PlayBack,
			d.MediaDuration, d.CIndex, d.DepId, d.Creater, d.Updater, d.Deleted)
		keys = append(keys, fmt.Sprintf("%s%v", cacheCourseCatalogueIdPrefix, d.Id))
	}
	query := fmt.Sprintf("insert into %s (%s) values %s", m.table, courseCatalogueRowsExpectAutoSet, strings.Join(values, ","))
	cc := m.WithSession(session)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		return conn.ExecCtx(ctx, query, args...)
	}, keys...)
	return err
}
