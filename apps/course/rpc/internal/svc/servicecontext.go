package svc

import (
	"context"
	"fmt"

	"tjxt/apps/course/rpc/internal/config"
	"tjxt/apps/course/rpc/internal/model"
	"tjxt/pkg/mq"
	"tjxt/pkg/mq/event"

	"github.com/zeromicro/go-zero/core/logx"
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

	// Producer 课程上下架事件发布者，可为 nil：RabbitMQ 未配置或连接失败时
	// 跳过发布（best-effort），不阻塞课程服务启动与课程主流程。
	// 发布的事件由 search 服务消费，增量同步 ES 课程索引。
	Producer *mq.Producer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)
	svcCtx := &ServiceContext{
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
	initProducer(svcCtx)
	return svcCtx
}

// initProducer 创建课程上架/下架事件发布者（course.events 交换机）。
// RabbitMQ 未配置或连接不可用时仅告警并置 nil，不阻塞课程服务启动。
func initProducer(svcCtx *ServiceContext) {
	c := svcCtx.Config.RabbitMQ
	if c.Host == "" || c.Port == 0 {
		logx.Infof("skip course event producer: rabbitmq not configured")
		return
	}
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.User, c.Pass, c.Host, c.Port)
	p, err := mq.NewProducer(dsn)
	if err != nil {
		logx.Errorf("init course event producer failed (events will not be published): %v", err)
		return
	}
	svcCtx.Producer = p
	logx.Infof("course event producer initialized, exchange=%s", mq.ExchangeCourse)
}

// PublishCourseEvent 发布课程上下架事件（best-effort，失败仅告警）。
// up=true 发布 course.up（上架/索引写入），up=false 发布 course.down（下架/索引删除）。
// 不返回错误：课程主流程（DB 写入）已成功，搜索索引同步失败不应回滚课程操作；
// 数据恢复兜底为 search 服务启动全量重建索引 + 手动 ReindexCourses RPC。
func (s *ServiceContext) PublishCourseEvent(ctx context.Context, courseID int64, up bool) {
	if s.Producer == nil {
		return
	}
	rk := mq.RoutingKeyCourseDown
	if up {
		rk = mq.RoutingKeyCourseUp
	}
	evt := event.CourseEvent{CourseID: courseID}
	if err := s.Producer.Publish(ctx, mq.ExchangeCourse, rk, evt); err != nil {
		logx.Errorf("publish course event failed, courseId=%d up=%v: %v", courseID, up, err)
	}
}