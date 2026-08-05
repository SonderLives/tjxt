# 公共库契约 (pkg/)

> 版本：v1.0 | 更新：2026-08-05

## pkg/auth - JWT 认证上下文

### 对外接口
```go
package auth

// GetUserId 从 context 获取用户 ID
func GetUserId(ctx context.Context) int64

// GetRole 从 context 获取角色编码
func GetRole(ctx context.Context) string

// GetPermissions 从 context 获取权限列表
func GetPermissions(ctx context.Context) []string

// SetUserContext 设置用户信息到 context（中间件用）
func SetUserContext(ctx context.Context, userId int64, role string, perms []string) context.Context
```

### 使用约定
- Handler/Logic 只能通过以上函数获取身份，**禁止**直接解析 JWT
- 中间件在 `internal/middleware/auth.go` 中实现 Token 验证并调用 `SetUserContext`

---

## pkg/response - 统一响应

### 结构定义
```go
type R struct {
    Code      int    `json:"code"`
    Msg       string `json:"msg"`
    RequestId string `json:"requestId"`
    Data      any    `json:"data"`
}
```

### 写入函数
```go
// Write 自动包装成 R 并写入 HTTP 响应
func Write(w http.ResponseWriter, r *http.Request, data any, err error)
```

### 行为
- `err == nil`：Code=0, Msg="success", Data=data
- `err == xerr.ErrCode`：Code=err.Code, Msg=err.Msg, Data=nil
- 其他 error：Code=100001, Msg=err.Error(), Data=nil
- 自动注入 `RequestId`（从 context 或 header 获取）

---

## pkg/utils/page - 分页契约

### 请求
```go
type PageRequest struct {
    PageNo   int `form:"pageNo,default=1"   validate:"gte=1"`
    PageSize int `form:"pageSize,default=20" validate:"gte=1,lte=100"`
}
```

### 响应
```go
type PageResponse[T any] struct {
    Total    int64 `json:"total"`
    List     []T   `json:"list"`
    PageNo   int   `json:"pageNo"`
    PageSize int   `json:"pageSize"`
}
```

### 使用
```go
// Logic 中
resp, total := logic.List(ctx, req.PageNo, req.PageSize)
return &page.PageResponse[Item]{Total: total, List: resp, PageNo: req.PageNo, PageSize: req.PageSize}, nil
```

---

## pkg/utils/idgen - 雪花算法 ID 生成

### 接口
```go
// NextId 生成唯一 ID (int64)
func NextId() int64

// NextIdStr 生成字符串 ID
func NextIdStr() string
```

### 使用场景
- 订单号、流水号、券码、文件 ID 等全局唯一标识
- **禁止**自行实现 ID 生成逻辑

---

## pkg/mq - RabbitMQ 事件总线

### 生产者接口
```go
type Producer interface {
    // Publish 发布事件（自动重试、确认机制）
    Publish(ctx context.Context, exchange, routingKey string, msg any) error
    
    // PublishDelay 延迟发布（基于死信队列）
    PublishDelay(ctx context.Context, exchange, routingKey string, msg any, delay time.Duration) error
}
```

### 消费者接口
```go
type Consumer interface {
    // Consume 注册消费者
    Consume(queue string, handler func(ctx context.Context, msg []byte) error) error
}
```

### 事件定义规范
```go
// 所有事件结构体放在 pkg/mq/event/
type OrderCreatedEvent struct {
    OrderId   int64  `json:"orderId"`
    UserId    int64  `json:"userId"`
    Amount    int64  `json:"amount"`
    Timestamp int64  `json:"timestamp"`
}
```

### 命名约定
- Exchange: `<domain>.events` (如 `trade.events`)
- RoutingKey: `<domain>.<action>` (如 `order.created`, `payment.paid`)
- Queue: `<svc>.<domain>.<action>` (如 `learning.order.created`)

---

## pkg/xerr - 统一错误码

见 [error-codes.md](error-codes.md)

---

## 版本兼容策略

| 变更类型 | 策略 |
|----------|------|
| 新增导出函数/类型 | 向下兼容，直接发布 |
| 修改函数签名 | 新增函数，旧函数标记 `// Deprecated` 保留 1 个大版本 |
| 删除导出符号 | 标记 Deprecated 1 个大版本后再删除 |
| 结构体新增字段 | 向下兼容（JSON 忽略未知字段） |
| 结构体删除字段 | 标记 Deprecated，保留字段但不再赋值 |

所有共享库变更需在 `pkg/CHANGELOG.md` 记录，发布时打 Tag `pkg/v<version>`。