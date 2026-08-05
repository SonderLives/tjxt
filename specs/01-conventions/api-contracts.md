# API 统一契约规范

> 版本：v1.0 | 更新：2026-08-05 | 适用：所有 `apps/*/api/*.api` 与生成的 handler

## 统一响应格式

所有 HTTP 接口统一返回 `pkg/response.R` 结构：

```go
// pkg/response/json.go
type R struct {
    Code      int    `json:"code"`       // 业务码：0=成功，非0=失败（见 error-codes.md）
    Msg       string `json:"msg"`        // 人类可读提示
    RequestId string `json:"requestId"`  // 链路追踪 ID
    Data      any    `json:"data"`       // 业务载荷，成功时非 nil
}
```

**Handler 必须使用**：
```go
// ✅ 正确
func (h *UserHandler) GetCurrent(w http.ResponseWriter, r *http.Request) {
    resp, err := h.logic.GetCurrent(r.Context())
    result.Write(w, r, resp, err)  // 自动包装成 R{Code,Msg,RequestId,Data}
}

// ❌ 禁止
httpx.OkJsonCtx(r.Context(), w, resp)        // 无统一码、无 requestId
httpx.ErrorCtx(r.Context(), w, err)          // 无统一错误格式
```

---

## 统一分页契约

### 请求参数（`pkg/utils/page.PageRequest`）
```go
type PageRequest struct {
    PageNo   int `form:"pageNo,default=1"   validate:"gte=1"`
    PageSize int `form:"pageSize,default=20" validate:"gte=1,lte=100"`
}
```

### 响应结构（`pkg/utils/page.PageResponse[T]`）
```go
type PageResponse[T any] struct {
    Total    int64 `json:"total"`
    List     []T   `json:"list"`
    PageNo   int   `json:"pageNo"`
    PageSize int   `json:"pageSize"`
}
```

### .api 定义示例
```go
type PageReq {
    PageNo   int `json:"pageNo,default=1"`
    PageSize int `json:"pageSize,default=20"`
}

type PageResp {
    Total    int64       `json:"total"`
    List     interface{} `json:"list"`  // 具体类型在 logic 里指定
    PageNo   int         `json:"pageNo"`
    PageSize int         `json:"pageSize"`
}
```

---

## 认证与授权

### JWT 保护
在 `.api` 文件的 `@server` 指令添加：
```go
@server(
    jwt: Auth
)
```

### 从 Context 取身份（`pkg/auth/context.go`）
```go
// ✅ 正确：logic 中
userId := auth.GetUserId(ctx)    // int64
role := auth.GetRole(ctx)        // string

// ❌ 禁止：handler 中手写 JWT 解析
token := r.Header.Get("Authorization")
// 手动 parse...
```

### 权限标签（RBAC）
- 接口级权限在 `.api` 的 `doc` 或自定义 annotation 标记
- 实际校验在 logic 或中间件，权限码格式：`<resource>:<action>`（如 `course:create`、`user:delete`）

---

## 错误码规范

| 分类 | 范围 | 说明 |
|------|------|------|
| 系统通用 | 100000-100999 | 参数校验、鉴权、限流、熔断 |
| 用户域 | 101000-101999 | 登录、注册、Token、角色 |
| 课程域 | 102000-102999 | 课程 CRUD、章节、资源 |
| 交易域 | 103000-103999 | 订单、支付、退款、分账 |
| 学习域 | 104000-104999 | 进度、笔记、证书 |
| 支付域 | 105000-105999 | 渠道、流水、对账 |
| 媒体域 | 106000-106999 | 上传、签名、转码 |
| 营销域 | 107000-107999 | 优惠券、活动 |
| 消息域 | 108000-108999 | 站内信、短信、模板 |
| 考试域 | 109000-109999 | 题库、试卷、考试 |
| 搜索域 | 110000-110999 | 兴趣、推荐 |
| 用户中心 | 111000-111999 | 档案、后台管理 |

**完整错误码表**：见 [03-shared/error-codes.md](../03-shared/error-codes.md)

---

## .api 文件编写规范

### 类型定义
```go
// 请求/响应类型命名：<动作><资源>Req/Resp
type CreateCourseReq { ... }
type CreateCourseResp { ... }

// 枚举用 string + comment
type CourseStatus string // NORMAL=正常, DRAFT=草稿, OFFLINE=下架

// goctl 不支持 any，用 interface{}
type CommonResp {
    Data interface{} `json:"data"`
}
```

### 路由分组与中间件
```go
@server(
    prefix: /api/v1
    jwt: Auth
    middleware: Logging,Recovery
)
```

### 路由定义
```go
// RESTful 风格
post  /courses           (CreateCourseReq) returns (CreateCourseResp)
get   /courses           (PageReq)         returns (PageResp)
get   /courses/:id       (GetCourseReq)    returns (CourseDetailResp)
put   /courses/:id       (UpdateCourseReq) returns (CourseDetailResp)
delete /courses/:id      (DeleteCourseReq) returns (CommonResp)

// 非 CRUD 动作用动词前缀
post  /courses/:id/publish     (PublishCourseReq) returns (CommonResp)
post  /courses/:id/offline     (OfflineCourseReq) returns (CommonResp)
```

### 文档注释（生成 Swagger 用）
```go
@doc(
    summary: "创建课程"
    description: "教师创建新课程，初始状态为草稿"
    tag: "课程管理"
)
post /courses (CreateCourseReq) returns (CreateCourseResp)
```

---

## 请求参数绑定规则

| 位置 | 绑定方式 | 示例 |
|------|----------|------|
| Path | `:id` → `Id int64` | `/courses/:id` |
| Query | `form:"pageNo"` | `?pageNo=2` |
| Header | `header:"X-Request-Id"` | `RequestId string` |
| Body | JSON 自动绑定 | `CreateCourseReq` |

---

## 版本与兼容

- URL 前缀固定 `/api/v1`，大版本升级新增 `/api/v2`
- 响应字段**只增不减**，废弃字段标记 `deprecated: true`
- 新增必填字段需提供默认值或兼容旧客户端