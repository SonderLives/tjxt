package model

import (
	"context"
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
