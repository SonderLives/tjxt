package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"strconv"
	"time"

	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/pkg/xerr"

	"github.com/elastic/go-elasticsearch/v9/esutil"
	"github.com/zeromicro/go-zero/core/logx"
)

const reindexPageSize = 500

// ReindexAll 从 course 服务分页拉取全部已上架课程，批量写入 ES course 索引。
// 返回成功写入的文档数与扫描到的课程总数。
// ES 不可用时返回 ServiceUnavailable。
// 幂等：以 courseId 作为文档 ID 做 upsert，重复执行安全（启动时 + 手动均可调用）。
func ReindexAll(ctx context.Context, svcCtx *ServiceContext) (indexed, total int64, err error) {
	if svcCtx.ES == nil {
		return 0, 0, xerr.ServiceUnavailable("搜索服务暂不可用")
	}

	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         CourseIndexName,
		Client:        svcCtx.ES,
		NumWorkers:    4,
		FlushBytes:    5 << 20, // 5MB
		FlushInterval: 10 * time.Second,
	})
	if err != nil {
		return 0, 0, xerr.Wrapf(err, xerr.CodeInternal, "初始化批量索引器失败")
	}

	var failed int64
	for pageNo := int64(1); ; pageNo++ {
		reply, err := svcCtx.CourseRpc.CourseSearchIndexInfoList(ctx, &courseclient.CourseSearchIndexInfoListRequest{
			PageNo:   pageNo,
			PageSize: reindexPageSize,
		})
		if err != nil {
			_ = bi.Close(ctx)
			return 0, total, xerr.Wrapf(err, xerr.CodeInternal, "拉取课程索引列表失败")
		}

		for _, info := range reply.Items {
			total++
			if info.Enable != 1 {
				continue
			}
			doc := courseDocFromInfo(info)
			body, mErr := json.Marshal(doc)
			if mErr != nil {
				failed++
				logx.Errorf("marshal course doc failed, courseId=%d: %v", doc.ID, mErr)
				continue
			}
			id := strconv.FormatInt(doc.ID, 10)
			if aErr := bi.Add(ctx, esutil.BulkIndexerItem{
				Action:     "index",
				DocumentID: id,
				Body:       bytes.NewReader(body),
			}); aErr != nil {
				_ = bi.Close(ctx)
				return 0, total, xerr.Wrapf(aErr, xerr.CodeInternal, "批量写入 ES 失败")
			}
		}

		if len(reply.Items) < reindexPageSize {
			break
		}
	}

	if err := bi.Close(ctx); err != nil {
		return 0, total, xerr.Wrapf(err, xerr.CodeInternal, "刷新批量索引器失败")
	}

	stats := bi.Stats()
	indexed = int64(stats.NumIndexed)
	failed += int64(stats.NumFailed)
	logx.Infof("reindex finished, total=%d indexed=%d failed=%d", total, indexed, failed)
	return indexed, total, nil
}

// initReindex 启动时在后台异步做一次全量重建索引（best-effort）：
// 让 course 索引在服务启动后即可被搜索，而不必依赖 RabbitMQ 的 course.up 事件。
// 失败仅告警，不阻塞服务启动；后续 MQ 事件仍会持续同步增量。
func initReindex(svcCtx *ServiceContext) {
	if svcCtx.ES == nil {
		logx.Error("skip startup reindex: elasticsearch unavailable")
		return
	}
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()
		indexed, total, err := ReindexAll(ctx, svcCtx)
		if err != nil {
			logx.Errorf("startup reindex failed: %v", err)
			return
		}
		logx.Infof("startup reindex done, total=%d indexed=%d", total, indexed)
	}()
}
