package svc

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/apps/search/rpc/internal/config"
	"tjxt/apps/search/rpc/internal/model"
	"tjxt/apps/search/rpc/pb"
	"tjxt/pkg/mq"

	"github.com/elastic/go-elasticsearch/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

// CourseIndexName ES 课程索引名。
const CourseIndexName = "course"

// CourseDoc ES 课程索引文档，字段与 course 服务 CourseSearchIndexInfo 1:1 对齐。
type CourseDoc struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	CoverURL      string `json:"cover_url"`
	Price         int64  `json:"price"`
	Score         int64  `json:"score"`
	Sold          int64  `json:"sold"`
	Sections      int64  `json:"sections"`
	Free          int32  `json:"free"`
	CourseType    int32  `json:"course_type"`
	Enable        int32  `json:"enable"`
	CategoryIDLv1 int64  `json:"category_id_lv1"`
	CategoryIDLv2 int64  `json:"category_id_lv2"`
	CategoryIDLv3 int64  `json:"category_id_lv3"`
	CreateTime    string `json:"create_time"`
	PublishTime   string `json:"publish_time"`
	Duration      int64  `json:"duration"`
}

// ToCourseVO 转换为 search RPC 的 CourseVO。
func (d *CourseDoc) ToCourseVO() *pb.CourseVO {
	return &pb.CourseVO{
		Id:            d.ID,
		Name:          d.Name,
		CoverUrl:      d.CoverURL,
		Price:         d.Price,
		Score:         d.Score,
		Sold:          d.Sold,
		Sections:      d.Sections,
		Free:          d.Free,
		CourseType:    d.CourseType,
		Enable:        d.Enable,
		CategoryIdLv1: d.CategoryIDLv1,
		CategoryIdLv2: d.CategoryIDLv2,
		CategoryIdLv3: d.CategoryIDLv3,
		CreateTime:    d.CreateTime,
		PublishTime:   d.PublishTime,
		Duration:      d.Duration,
	}
}

// courseIndexMapping 课程索引 mapping。
// name 字段使用配置的 analyzer（standard：中文按单字切分，探测结论见
// etc/search.yaml），并挂 keyword 子字段支持精确匹配/排序。
// create_time/publish_time 与 course 服务返回的字符串格式一致，映射为
// date（yyyy-MM-dd HH:mm:ss）便于范围查询。
const courseIndexMapping = `{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0
  },
  "mappings": {
    "properties": {
      "id": {"type": "long"},
      "name": {
        "type": "text",
        "analyzer": "%s",
        "fields": {"keyword": {"type": "keyword", "ignore_above": 256}}
      },
      "cover_url": {"type": "keyword", "ignore_above": 1024},
      "price": {"type": "long"},
      "score": {"type": "long"},
      "sold": {"type": "long"},
      "sections": {"type": "long"},
      "free": {"type": "long"},
      "course_type": {"type": "long"},
      "enable": {"type": "long"},
      "category_id_lv1": {"type": "long"},
      "category_id_lv2": {"type": "long"},
      "category_id_lv3": {"type": "long"},
      "create_time": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
      "publish_time": {"type": "date", "format": "yyyy-MM-dd HH:mm:ss"},
      "duration": {"type": "long"}
    }
  }
}`

type ServiceContext struct {
	Config         config.Config
	InterestsModel model.InterestsModel

	// CourseRpc 课程服务客户端，索引同步时回源课程数据
	CourseRpc courseclient.Course

	// ES ES 客户端，可为 nil：初始化失败时检索接口返回 503，
	// 且跳过 MQ 索引同步（不阻塞服务启动）
	ES *elasticsearch.Client

	// MQClient 消费 course.events 的客户端，可为 nil：ES 不可用时跳过
	MQClient *mq.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.DataSource)

	svcCtx := &ServiceContext{
		Config:         c,
		InterestsModel: model.NewInterestsModel(conn, c.Cache),
		CourseRpc:      courseclient.NewCourse(zrpc.MustNewClient(c.CourseRpc)),
	}

	initES(svcCtx)
	initMQ(svcCtx)

	return svcCtx
}

// initES 初始化 ES 客户端：连通性自检 + 幂等创建 course 索引。
// 失败仅告警，不阻塞服务启动。
func initES(svcCtx *ServiceContext) {
	cfg := svcCtx.Config.Elasticsearch
	analyzer := strings.TrimSpace(cfg.Analyzer)
	if analyzer == "" {
		analyzer = "standard"
	}

	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
		Username:  cfg.Username,
		Password:  cfg.Password,
	})
	if err != nil {
		logx.Errorf("init elasticsearch client failed: %v", err)
		return
	}

	// 连通性自检
	info, err := es.Info()
	if err != nil {
		logx.Errorf("elasticsearch self-check failed: %v", err)
		return
	}
	infoBody, _ := io.ReadAll(info.Body)
	info.Body.Close()

	svcCtx.ES = es
	logx.Infof("elasticsearch connected: %s, analyzer=%s", strings.TrimSpace(string(infoBody)), analyzer)

	if err := ensureCourseIndex(es, analyzer); err != nil {
		logx.Errorf("ensure elasticsearch index %q failed: %v", CourseIndexName, err)
	}
}

// ensureCourseIndex 幂等创建 course 索引（Exists → Create）。
func ensureCourseIndex(es *elasticsearch.Client, analyzer string) error {
	exists, err := es.Indices.Exists([]string{CourseIndexName})
	if err != nil {
		return fmt.Errorf("check index exists: %w", err)
	}
	exists.Body.Close()
	if exists.StatusCode == http.StatusOK {
		logx.Infof("elasticsearch index %q already exists", CourseIndexName)
		return nil
	}

	mapping := fmt.Sprintf(courseIndexMapping, analyzer)
	res, err := es.Indices.Create(CourseIndexName, es.Indices.Create.WithBody(bytes.NewReader([]byte(mapping))))
	if err != nil {
		return fmt.Errorf("create index: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("create index %q: %s", CourseIndexName, res.String())
	}
	logx.Infof("elasticsearch index %q created, analyzer=%s", CourseIndexName, analyzer)
	return nil
}

// initMQ 注册课程上下架事件消费者。ES 不可用时跳过注册。
func initMQ(svcCtx *ServiceContext) {
	if svcCtx.ES == nil {
		logx.Error("skip mq course event consumer: elasticsearch unavailable")
		return
	}

	c := svcCtx.Config
	dsn := fmt.Sprintf("amqp://%s:%s@%s:%d/", c.RabbitMQ.User, c.RabbitMQ.Pass, c.RabbitMQ.Host, c.RabbitMQ.Port)
	client := mq.NewClient(dsn)

	consumer := &courseIndexConsumer{svcCtx: svcCtx}
	mq.Register(client, mq.Binding{
		Queue:      "search.course.up",
		Exchange:   mq.ExchangeCourse,
		RoutingKey: mq.RoutingKeyCourseUp,
	}, consumer.handleCourseUp)
	mq.Register(client, mq.Binding{
		Queue:      "search.course.down",
		Exchange:   mq.ExchangeCourse,
		RoutingKey: mq.RoutingKeyCourseDown,
	}, consumer.handleCourseDown)

	svcCtx.MQClient = client
}
