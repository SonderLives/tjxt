package service

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"tjxt/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

// CourseSimpleInfo 课程简要信息（来自课程服务内部接口 CourseSimpleInfoDTO）
type CourseSimpleInfo struct {
	Id            int64  `json:"id"`
	Name          string `json:"name"`
	CoverUrl      string `json:"coverUrl"`
	Price         int64  `json:"price"`
	Free          bool   `json:"free"`
	ValidDuration int64  `json:"validDuration"` // 有效期，单位月，0 表示永久
	Status        int64  `json:"status"`
}

// CourseClient 课程服务客户端抽象
type CourseClient interface {
	GetSimpleInfos(ctx context.Context, ids []int64) (map[int64]*CourseSimpleInfo, error)
}

// UserInfo 用户简要信息（来自用户服务内部接口 UserDTO）
type UserInfo struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CellPhone string `json:"cellPhone"`
	Type      int64  `json:"type"`
}

// UserClient 用户服务客户端抽象
type UserClient interface {
	GetByIds(ctx context.Context, ids []int64) (map[int64]*UserInfo, error)
}

// httpClient 内部服务 HTTP 客户端
type httpClient struct {
	baseURL string
	timeout time.Duration
	client  *http.Client
}

func newHTTPClient(baseURL string, timeoutMillis int64) *httpClient {
	if timeoutMillis <= 0 {
		timeoutMillis = 3000
	}
	return &httpClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		timeout: time.Duration(timeoutMillis) * time.Millisecond,
		client:  &http.Client{Timeout: time.Duration(timeoutMillis) * time.Millisecond},
	}
}

// apiResponse 内部服务的统一响应结构
type apiResponse struct {
	Code      int64           `json:"code"`
	Msg       string          `json:"msg"`
	RequestId string          `json:"requestId"`
	Data      json.RawMessage `json:"data"`
}

// getJSON 发起 GET 请求并解析统一响应，目标 schema 失败返回 503。
func (h *httpClient) getJSON(ctx context.Context, path string, query url.Values, target any) error {
	reqURL := h.baseURL + path
	if len(query) > 0 {
		reqURL += "?" + query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, reqURL, nil)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "构建内部请求失败")
	}
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return xerr.Wrap(err, xerr.CodeServiceUnavailable, "依赖服务暂不可用")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "读取内部服务响应失败")
	}

	var r apiResponse
	if err := json.Unmarshal(body, &r); err != nil {
		return xerr.Wrap(err, xerr.CodeInternal, "解析内部服务响应失败")
	}
	if r.Code != int64(xerr.CodeSuccess) {
		logx.Errorf("internal service error, url=%s, code=%d, msg=%s", reqURL, r.Code, r.Msg)
		return xerr.ServiceUnavailable("依赖服务返回错误")
	}
	if target != nil && len(r.Data) > 0 {
		if err := json.Unmarshal(r.Data, target); err != nil {
			return xerr.Wrap(err, xerr.CodeInternal, "解析内部服务数据失败")
		}
	}
	return nil
}

// courseHTTPClient 基于课程服务内部 HTTP 接口的实现
type courseHTTPClient struct {
	http *httpClient
}

// NewCourseHTTPClient 创建课程服务 HTTP 客户端。
func NewCourseHTTPClient(baseURL string, timeoutMillis int64) CourseClient {
	return &courseHTTPClient{http: newHTTPClient(baseURL, timeoutMillis)}
}

// GetSimpleInfos 批量查询课程简要信息。
func (c *courseHTTPClient) GetSimpleInfos(ctx context.Context, ids []int64) (map[int64]*CourseSimpleInfo, error) {
	result := make(map[int64]*CourseSimpleInfo, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	var list []CourseSimpleInfo
	q := url.Values{}
	q.Set("ids", joinIDs(ids))
	if err := c.http.getJSON(ctx, "/courses/simpleInfo/list", q, &list); err != nil {
		return nil, err
	}
	for i := range list {
		result[list[i].Id] = &list[i]
	}
	return result, nil
}

// userHTTPClient 基于用户服务内部 HTTP 接口的实现
type userHTTPClient struct {
	http *httpClient
}

// NewUserHTTPClient 创建用户服务 HTTP 客户端。
func NewUserHTTPClient(baseURL string, timeoutMillis int64) UserClient {
	return &userHTTPClient{http: newHTTPClient(baseURL, timeoutMillis)}
}

// GetByIds 批量查询用户简要信息。
func (c *userHTTPClient) GetByIds(ctx context.Context, ids []int64) (map[int64]*UserInfo, error) {
	result := make(map[int64]*UserInfo, len(ids))
	if len(ids) == 0 {
		return result, nil
	}
	for _, id := range ids {
		if _, ok := result[id]; ok {
			continue
		}
		var u UserInfo
		path := "/users/" + strconv.FormatInt(id, 10)
		if err := c.http.getJSON(ctx, path, nil, &u); err != nil {
			logx.Errorf("fetch user %d failed: %v", id, err)
			continue
		}
		u.Id = id
		result[id] = &u
	}
	return result, nil
}

func joinIDs(ids []int64) string {
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		parts = append(parts, strconv.FormatInt(id, 10))
	}
	return strings.Join(parts, ",")
}
