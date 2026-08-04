package svc

import (
	"tjxt/apps/course/rpc/internal/config"
	"tjxt/apps/course/rpc/internal/model"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	CategoryModel              model.CategoryModel
	CourseModel                model.CourseModel
	CourseDraftModel           model.CourseDraftModel
	CourseContentModel         model.CourseContentModel
	CourseContentDraftModel    model.CourseContentDraftModel
	CourseCatalogueModel       model.CourseCatalogueModel
	CourseCatalogueDraftModel  model.CourseCatalogueDraftModel
	CourseCataSubjectDraftModel model.CourseCataSubjectDraftModel
	CourseSubjectModel         model.CourseSubjectModel
	CourseTeacherModel         model.CourseTeacherModel
	CourseTeacherDraftModel    model.CourseTeacherDraftModel
	SubjectModel               model.SubjectModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	return &ServiceContext{
		Config:                    c,
		CategoryModel:             model.NewCategoryModel(conn, c.Cache),
		CourseModel:               model.NewCourseModel(conn, c.Cache),
		CourseDraftModel:          model.NewCourseDraftModel(conn, c.Cache),
		CourseContentModel:        model.NewCourseContentModel(conn, c.Cache),
		CourseContentDraftModel:   model.NewCourseContentDraftModel(conn, c.Cache),
		CourseCatalogueModel:      model.NewCourseCatalogueModel(conn, c.Cache),
		CourseCatalogueDraftModel: model.NewCourseCatalogueDraftModel(conn, c.Cache),
		CourseCataSubjectDraftModel: model.NewCourseCataSubjectDraftModel(conn, c.Cache),
		CourseSubjectModel:        model.NewCourseSubjectModel(conn, c.Cache),
		CourseTeacherModel:        model.NewCourseTeacherModel(conn, c.Cache),
		CourseTeacherDraftModel:   model.NewCourseTeacherDraftModel(conn, c.Cache),
		SubjectModel:              model.NewSubjectModel(conn, c.Cache),
	}
}