package svc

import (
	"course/internal/config"
	"course/internal/model"
	"course/internal/service"

	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 课程服务上下文：数据访问对象与业务服务统一装配。
type ServiceContext struct {
	Config config.Config

	// 数据访问
	CategoryModel          *model.CategoryModel
	CourseModel            *model.CourseModel
	CourseDraftModel       *model.CourseDraftModel
	CourseCatalogueModel   *model.CourseCatalogueModel
	CourseCatalogueDraftModel *model.CourseCatalogueDraftModel
	CourseTeacherModel     *model.CourseTeacherModel
	CourseTeacherDraftModel    *model.CourseTeacherDraftModel
	CourseContentModel     *model.CourseContentModel
	CourseContentDraftModel    *model.CourseContentDraftModel
	CataSubjectDraftModel  *model.CourseCataSubjectDraftModel
	SubjectModel           *model.SubjectModel
	UserDetailModel        *model.UserDetailModel

	// 业务服务
	CategoryService   service.CategoryService
	CourseService     service.CourseService
	CatalogueService  service.CatalogueService
}

// NewServiceContext 创建服务上下文。
func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)

	categoryModel := model.NewCategoryModel(conn)
	courseModel := model.NewCourseModel(conn)
	courseDraftModel := model.NewCourseDraftModel(conn)
	cataModel := model.NewCourseCatalogueModel(conn)
	cataDraftModel := model.NewCourseCatalogueDraftModel(conn)
	teacherModel := model.NewCourseTeacherModel(conn)
	teacherDraftModel := model.NewCourseTeacherDraftModel(conn)
	contentModel := model.NewCourseContentModel(conn)
	contentDraftModel := model.NewCourseContentDraftModel(conn)
	cataSubjectDraftModel := model.NewCourseCataSubjectDraftModel(conn)
	subjectModel := model.NewSubjectModel(conn)
	userDetailModel := model.NewUserDetailModel(conn, c.UserDetailTable)

	trailerDuration := c.TrailerDuration
	if trailerDuration <= 0 {
		trailerDuration = 5
	}

	return &ServiceContext{
		Config:                  c,
		CategoryModel:           categoryModel,
		CourseModel:             courseModel,
		CourseDraftModel:        courseDraftModel,
		CourseCatalogueModel:    cataModel,
		CourseCatalogueDraftModel: cataDraftModel,
		CourseTeacherModel:      teacherModel,
		CourseTeacherDraftModel: teacherDraftModel,
		CourseContentModel:      contentModel,
		CourseContentDraftModel: contentDraftModel,
		CataSubjectDraftModel:   cataSubjectDraftModel,
		SubjectModel:            subjectModel,
		UserDetailModel:         userDetailModel,
		CategoryService: service.NewCategoryService(categoryModel, courseModel, courseDraftModel),
		CourseService: service.NewCourseService(
			courseDraftModel, courseModel, contentDraftModel, contentModel,
			cataDraftModel, cataModel, teacherDraftModel, teacherModel,
			cataSubjectDraftModel, categoryModel, subjectModel, userDetailModel, trailerDuration),
		CatalogueService: service.NewCatalogueService(cataModel, trailerDuration),
	}
}
