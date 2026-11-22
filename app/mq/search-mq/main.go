package main

import (
	"context"
	"flag"
	"fmt"

	"tjxt/api/search/backend/svc"
	"tjxt/app/mq/search-mq/internal/config"

	"github.com/zeromicro/go-zero/core/logx"
	"tjxt/pkg/confenv"
)

var configFile = flag.String("f", "etc/search-mq.yaml", "the config file")

// search-mq 课程索引同步消费者（原 search rpc 进程内的消费者部分独立成进程）：
// 订阅 course.events 的 course.up / course.down 事件写/删 ES 索引，
// 启动时后台做一次全量重建（见 backend/svc 的 initMQ/initReindex）。
func main() {
	flag.Parse()

	var c config.Config
	confenv.MustLoad(*configFile, &c)
	// 无 rest/zrpc 服务器的纯消费者进程，需自行初始化 logx/prometheus/telemetry
	c.MustSetUp()

	ctx := svc.NewServiceContext(c.Config)

	// ES 不可用时 MQClient 为 nil（svc 初始化已告警），消费者无职责可履行，直接退出
	if ctx.MQClient == nil {
		logx.Error("search-mq exits: mq client unavailable (elasticsearch down?)")
		return
	}

	fmt.Printf("Starting search-mq consumer...\n")
	// 阻塞消费；连接断开后内部自动重连，进程退出依赖信号
	if err := ctx.MQClient.Start(context.Background()); err != nil {
		logx.Errorf("search mq consumer stopped: %v", err)
	}
}
