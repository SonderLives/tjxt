package svc

import (
	"course/internal/config"
	"course/internal/model"
	"course/internal/service"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

// ServiceContext 课程服务上下文：数据访问对象与业务服务统一装配。
type ServiceContext struct {
	Config config.Config

	// 数据访问
	CategoryModel            model.CategoryModel
	CourseModel              model.CourseModel
	CourseDraftModel         model.CourseDraftModel
	CourseCatalogueModel     model.CourseCatalogueModel
	CourseCatalogueDraftModel model.CourseCatalogueDraftModel
	CourseTeacherModel       model.CourseTeacherModel
	CourseTeacherDraftModel  model.CourseTeacherDraftModel
	CourseContentModel       model.CourseContentModel
	CourseContentDraftModel  model.CourseContentDraftModel
	CataSubjectDraftModel    model.CourseCataSubjectDraftModel
	SubjectModel             model.SubjectModel
	UserDetailModel          *model.UserDetailModel

	// 业务服务
	CategoryService  service.CategoryService
	CourseService    service.CourseService
	CatalogueService service.CatalogueService
}

// NewServiceContext 创建服务上下文。
func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DB.DataSource)

	// Redis 缓存配置
	cacheConf := cache.CacheConf{
		Host: c.Redis.Host,
		Type: c.Redis.Type,
		Pass: c.Redis.Pass,
	}

	categoryModel := model.NewCategoryModel(conn, cacheConf)
	courseModel := model.NewCourseModel(conn, cacheConf)
	courseDraftModel := model.NewCourseDraftModel(conn, cacheConf)
	cataModel := model.NewCourseCatalogueModel(conn, cacheConf)
	cataDraftModel := model.NewCourseCatalogueDraftModel(conn, cacheConf)
	teacherModel := model.NewCourseTeacherModel(conn, cacheConf)
	teacherDraftModel := model.NewCourseTeacherDraftModel(conn, cacheConf)
	contentModel := model.NewCourseContentModel(conn, cacheConf)
	contentDraftModel := model.NewCourseContentDraftModel(conn, cacheConf)
	cataSubjectDraftModel := model.NewCourseCataSubjectDraftModel(conn, cacheConf)
	subjectModel := model.NewSubjectModel(conn, cacheConf)
	userDetailModel := model.NewUserDetailModel(conn, c.UserDetailTable)

	trailerDuration := c.TrailerDuration
	if trailerDuration <= 0 {
		trailerDuration = 5
	}

	return &ServiceContext{
		Config:                     c,
		CategoryModel:              categoryModel,
		CourseModel:                courseModel,
		CourseDraftModel:           courseDraftModel,
		CourseCatalogueModel:       cataModel,
		CourseCatalogueDraftModel:  cataDraftModel,
		CourseTeacherModel:         teacherModel,
		CourseTeacherDraftModel:    teacherDraftModel,
		CourseContentModel:         contentModel,
		CourseContentDraftModel:    contentDraftModel,
		CataSubjectDraftModel:      cataSubjectDraftModel,
		SubjectModel:               subjectModel,
		UserDetailModel:            userDetailModel,
		CategoryService:            service.NewCategoryService(categoryModel, courseModel, courseDraftModel),
		CourseService:              service.NewCourseService(courseDraftModel, courseModel, contentDraftModel, contentModel, cataDraftModel, cataModel, teacherDraftModel, teacherModel, cataSubjectDraftModel, categoryModel, subjectModel, userDetailModel, trailerDuration),
		CatalogueService:           service.NewCatalogueService(cataModel, trailerDuration),
	}
}
