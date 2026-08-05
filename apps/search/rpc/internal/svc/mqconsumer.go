package svc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	courseclient "tjxt/apps/course/rpc/course"
	"tjxt/pkg/mq/event"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// 单条事件索引同步的最大重试次数，超过后丢弃消息并告警
const maxSyncRetries = 3

// courseIndexConsumer 课程索引同步消费者。
// 消费 course.events 交换机上的 course.up / course.down 事件，
// 将课程上架/下架状态同步到 ES 索引。
// 消费幂等：同一 courseId 重复事件只会执行一次 upsert/delete。
type courseIndexConsumer struct {
	svcCtx *ServiceContext
}

// handleCourseUp 课程上架：回源 course 服务索引数据 → ES upsert。
func (c *courseIndexConsumer) handleCourseUp(ctx context.Context, evt *event.CourseEvent) error {
	if evt.CourseID <= 0 {
		logx.Errorf("ignore course.up event with invalid courseId: %+v", evt)
		return nil
	}
	var lastErr error
	for i := 1; i <= maxSyncRetries; i++ {
		lastErr = c.syncUp(ctx, evt.CourseID)
		if lastErr == nil {
			return nil
		}
		logx.Errorf("sync course up to index failed, courseId=%d attempt=%d: %v", evt.CourseID, i, lastErr)
		time.Sleep(time.Duration(i) * time.Second)
	}
	logx.Errorf("discard course.up event after %d attempts, courseId=%d: %v", maxSyncRetries, evt.CourseID, lastErr)
	return nil
}

// handleCourseDown 课程下架：ES delete，文档不存在也幂等成功。
func (c *courseIndexConsumer) handleCourseDown(ctx context.Context, evt *event.CourseEvent) error {
	if evt.CourseID <= 0 {
		logx.Errorf("ignore course.down event with invalid courseId: %+v", evt)
		return nil
	}
	var lastErr error
	for i := 1; i <= maxSyncRetries; i++ {
		lastErr = c.syncDown(ctx, evt.CourseID)
		if lastErr == nil {
			return nil
		}
		logx.Errorf("sync course down from index failed, courseId=%d attempt=%d: %v", evt.CourseID, i, lastErr)
		time.Sleep(time.Duration(i) * time.Second)
	}
	logx.Errorf("discard course.down event after %d attempts, courseId=%d: %v", maxSyncRetries, evt.CourseID, lastErr)
	return nil
}

// syncUp 查询课程索引数据并写入 ES（upsert）。
func (c *courseIndexConsumer) syncUp(ctx context.Context, courseId int64) error {
	info, err := c.svcCtx.CourseRpc.CourseSearchInfoForIndex(ctx, &courseclient.IdRequest{Id: courseId})
	if err != nil {
		// 课程不存在视为终态：不上索引，直接丢弃（避免无效重试）
		if xerr.CodeOf(err) == xerr.CodeNotFound {
			logx.Errorf("course not found for index, courseId=%d: %v", courseId, err)
			return nil
		}
		return fmt.Errorf("query course index info: %w", err)
	}
	// 校验 enable==1，未上架课程不写入索引
	if info == nil || info.Enable != 1 {
		logx.Infof("skip indexing course with enable!=1, courseId=%d", courseId)
		return nil
	}

	doc := courseDocFromInfo(info)
	body, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("marshal course doc: %w", err)
	}

	res, err := c.svcCtx.ES.Index(CourseIndexName, bytes.NewReader(body),
		c.svcCtx.ES.Index.WithDocumentID(strconv.FormatInt(courseId, 10)))
	if err != nil {
		return fmt.Errorf("index course doc: %w", err)
	}
	defer res.Body.Close()
	if res.IsError() {
		return fmt.Errorf("index course doc: %s", res.String())
	}
	logx.Infof("course index upserted, courseId=%d", courseId)
	return nil
}

// syncDown 从 ES 删除课程文档，404 视为幂等成功。
func (c *courseIndexConsumer) syncDown(ctx context.Context, courseId int64) error {
	res, err := c.svcCtx.ES.Delete(CourseIndexName, strconv.FormatInt(courseId, 10))
	if err != nil {
		return fmt.Errorf("delete course doc: %w", err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return nil
	}
	if res.IsError() {
		return fmt.Errorf("delete course doc: %s", res.String())
	}
	logx.Infof("course index deleted, courseId=%d", courseId)
	return nil
}

// courseDocFromInfo 将 course 服务索引数据映射为 ES 文档（字段 1:1）。
func courseDocFromInfo(info *courseclient.CourseSearchIndexInfo) *CourseDoc {
	return &CourseDoc{
		ID:            info.Id,
		Name:          info.Name,
		CoverURL:      info.CoverUrl,
		Price:         info.Price,
		Score:         info.Score,
		Sold:          info.Sold,
		Sections:      info.Sections,
		Free:          info.Free,
		CourseType:    info.CourseType,
		Enable:        info.Enable,
		CategoryIDLv1: info.CategoryIdLv1,
		CategoryIDLv2: info.CategoryIdLv2,
		CategoryIDLv3: info.CategoryIdLv3,
		CreateTime:    info.CreateTime,
		PublishTime:   info.PublishTime,
		Duration:      info.Duration,
	}
}
