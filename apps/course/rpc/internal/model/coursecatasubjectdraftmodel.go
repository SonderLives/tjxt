package model

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ CourseCataSubjectDraftModel = (*customCourseCataSubjectDraftModel)(nil)

type (
	// CourseCataSubjectDraftModel is an interface to be customized, add more methods here,
	// and implement the added methods in customCourseCataSubjectDraftModel.
	CourseCataSubjectDraftModel interface {
		courseCataSubjectDraftModel
		ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCataSubjectDraft, error)
		ListByCataId(ctx context.Context, cataId int64) ([]*CourseCataSubjectDraft, error)
		DeleteByCourseId(ctx context.Context, courseId int64) error
		DeleteByCataId(ctx context.Context, cataId int64) error
	}
	customCourseCataSubjectDraftModel struct {
		*defaultCourseCataSubjectDraftModel
	}
)

// NewCourseCataSubjectDraftModel returns a model for the database table.
func NewCourseCataSubjectDraftModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) CourseCataSubjectDraftModel {
	return &customCourseCataSubjectDraftModel{
		defaultCourseCataSubjectDraftModel: newCourseCataSubjectDraftModel(conn, c, opts...),
	}
}

func (m *customCourseCataSubjectDraftModel) ListByCourseId(ctx context.Context, courseId int64) ([]*CourseCataSubjectDraft, error) {
	var list []*CourseCataSubjectDraft
	query := fmt.Sprintf("select * from %s where `course_id` = ?", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, courseId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCataSubjectDraftModel) ListByCataId(ctx context.Context, cataId int64) ([]*CourseCataSubjectDraft, error) {
	var list []*CourseCataSubjectDraft
	query := fmt.Sprintf("select * from %s where `cata_id` = ?", m.table)
	if err := m.QueryRowsNoCacheCtx(ctx, &list, query, cataId); err != nil {
		return nil, err
	}
	return list, nil
}

func (m *customCourseCataSubjectDraftModel) DeleteByCourseId(ctx context.Context, courseId int64) error {
	query := fmt.Sprintf("delete from %s where `course_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, courseId)
	return err
}

func (m *customCourseCataSubjectDraftModel) DeleteByCataId(ctx context.Context, cataId int64) error {
	query := fmt.Sprintf("delete from %s where `cata_id` = ?", m.table)
	_, err := m.ExecNoCacheCtx(ctx, query, cataId)
	return err
}
