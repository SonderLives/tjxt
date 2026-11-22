package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseContentModel = (*customCourseContentModel)(nil)

type (
	// CourseContentModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseContentModel.
	CourseContentModel interface {
		courseContentModel
		Upsert(ctx context.Context, data *CourseContent) error
		UpsertTx(ctx context.Context, session sqlx.Session, data *CourseContent) error
	}
	customCourseContentModel struct {
		*defaultCourseContentModel
	}
)

// NewCourseContentModel returns a model for the database table.
func NewCourseContentModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseContentModel {
	return &customCourseContentModel{
		defaultCourseContentModel: newCourseContentModel(conn, c, opts...),
	}
}

// Upsert 课程内容（大文本）按课程 id 主键替换插入或更新。
func (m *customCourseContentModel) Upsert(ctx context.Context, data *CourseContent) error {
	query := fmt.Sprintf("replace into %s (`id`,`course_introduce`,`use_people`,`course_detail`,`dep_id`,`creater`,`updater`,`deleted`,`create_time`,`update_time`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query,
		data.Id, data.CourseIntroduce, data.UsePeople, data.CourseDetail,
		data.DepId, data.Creater, data.Updater, data.Deleted, data.CreateTime, data.UpdateTime)
	return err
}

// UpsertTx 事务内按课程 id 替换插入课程内容，SQL 同 Upsert；写入走传入的事务会话，
// 并失效内容行缓存（供跨模型事务如课程上架发布复用）。
func (m *customCourseContentModel) UpsertTx(ctx context.Context, session sqlx.Session, data *CourseContent) error {
	cc := m.WithSession(session)
	courseContentIdKey := fmt.Sprintf("%s%v", cacheCourseContentIdPrefix, data.Id)
	_, err := cc.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("replace into %s (`id`,`course_introduce`,`use_people`,`course_detail`,`dep_id`,`creater`,`updater`,`deleted`,`create_time`,`update_time`) values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", m.table)
		return conn.ExecCtx(ctx, query,
			data.Id, data.CourseIntroduce, data.UsePeople, data.CourseDetail,
			data.DepId, data.Creater, data.Updater, data.Deleted, data.CreateTime, data.UpdateTime)
	}, courseContentIdKey)
	return err
}
