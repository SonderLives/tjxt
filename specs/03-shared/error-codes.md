# 共享错误码总表

> 版本：v1.0 | 更新：2026-08-05 | 来源：`pkg/xerr/code.go`

## 系统通用错误码 (100000-100999)

| 错误码 | 常量名 | HTTP Status | 说明 |
|--------|--------|-------------|------|
| 100000 | `SUCCESS` | 200 | 成功 |
| 100001 | `SERVER_ERROR` | 500 | 服务器内部错误 |
| 100002 | `PARAM_ERROR` | 400 | 参数校验失败 |
| 100003 | `UNAUTHORIZED` | 401 | 未登录/Token 失效 |
| 100004 | `FORBIDDEN` | 403 | 无权限访问 |
| 100005 | `NOT_FOUND` | 404 | 资源不存在 |
| 100006 | `METHOD_NOT_ALLOWED` | 405 | 请求方法不支持 |
| 100007 | `RATE_LIMITED` | 429 | 请求过于频繁 |
| 100008 | `SERVICE_UNAVAILABLE` | 503 | 服务不可用（熔断/降级） |
| 100009 | `DB_ERROR` | 500 | 数据库操作失败 |
| 100010 | `CACHE_ERROR` | 500 | 缓存操作失败 |
| 100011 | `MQ_ERROR` | 500 | 消息队列操作失败 |
| 100012 | `RPC_ERROR` | 500 | RPC 调用失败 |
| 100013 | `TIMEOUT` | 504 | 调用超时 |

## 业务错误码分段

| 领域 | 范围 | 文件 |
|------|------|------|
| 认证授权 | 101000-101999 | `pkg/xerr/auth.go` |
| 课程管理 | 102000-102999 | `pkg/xerr/course.go` |
| 交易订单 | 103000-103999 | `pkg/xerr/trade.go` |
| 学习中心 | 104000-104999 | `pkg/xerr/learning.go` |
| 支付网关 | 105000-105999 | `pkg/xerr/pay.go` |
| 媒体文件 | 106000-106999 | `pkg/xerr/media.go` |
| 营销优惠 | 107000-107999 | `pkg/xerr/promotion.go` |
| 消息通知 | 108000-108999 | `pkg/xerr/message.go` |
| 考试中心 | 109000-109999 | `pkg/xerr/exam.go` |
| 搜索推荐 | 110000-110999 | `pkg/xerr/search.go` |
| 用户中心 | 111000-111999 | `pkg/xerr/user.go` |

## 使用示例

```go
// logic 中返回错误
if user == nil {
    return nil, xerr.NewErrCode(xerr.USER_NOT_FOUND)
}

// 包装底层错误
if err != nil {
    return nil, xerr.NewErrCode(xerr.DB_ERROR).WithError(err)
}

// 自定义消息
return nil, xerr.NewErrMsg("自定义错误提示")
```

## 扩展规则

1. **新增错误码**：在对应领域文件中添加常量，保持范围不重叠
2. **错误信息**：Code 对应的 Msg 统一在 `xerr` 包内维护，logic 只需返回 Code
3. **国际化**：后续扩展支持多语言 Msg Map