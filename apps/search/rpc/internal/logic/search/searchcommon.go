package searchlogic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"tjxt/apps/search/rpc/internal/svc"
	"tjxt/apps/search/rpc/pb"
	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// courseSearch 执行 ES 查询并解析为分页结果。
// ES 调用失败统一返回 xerr.ServiceUnavailable。
func courseSearch(ctx context.Context, svcCtx *svc.ServiceContext, body []byte) (*pb.CourseSearchPageReply, error) {
	res, err := svcCtx.ES.Search(
		svcCtx.ES.Search.WithContext(ctx),
		svcCtx.ES.Search.WithIndex(svc.CourseIndexName),
		svcCtx.ES.Search.WithBody(bytes.NewReader(body)),
		svcCtx.ES.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeServiceUnavailable, "搜索服务暂不可用")
	}
	defer res.Body.Close()
	if res.IsError() {
		return nil, xerr.Wrapf(fmt.Errorf("%s", res.String()), xerr.CodeServiceUnavailable, "搜索服务暂不可用")
	}

	// 9.x 响应结构与 8.x 一致：total 为 {"value": n, "relation": "eq"}
	var parsed struct {
		Hits struct {
			Total struct {
				Value int64 `json:"value"`
			} `json:"total"`
			Hits []struct {
				Source json.RawMessage `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return nil, xerr.Wrapf(err, xerr.CodeServiceUnavailable, "搜索服务暂不可用")
	}

	reply := &pb.CourseSearchPageReply{Total: parsed.Hits.Total.Value}
	for _, h := range parsed.Hits.Hits {
		var doc svc.CourseDoc
		if err := json.Unmarshal(h.Source, &doc); err != nil {
			logx.Errorf("unmarshal course doc failed: %v", err)
			continue
		}
		reply.Items = append(reply.Items, doc.ToCourseVO())
	}
	return reply, nil
}
