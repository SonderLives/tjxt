# error-codes（共享错误码）Specification

## Purpose

定义全系统统一的业务错误码体系（`pkg/xerr`），所有服务的 logic 层统一返回错误码、由响应层统一映射 HTTP Status，保证错误语义跨服务一致。

## Requirements

### Requirement: 系统通用错误码

系统 SHALL 提供 100000-100999 范围的系统通用错误码，语义与 HTTP Status 对应关系如下：

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

#### Scenario: 资源不存在时返回 NOT_FOUND

- **WHEN** logic 层查询的资源不存在并以 `xerr.NOT_FOUND` 返回
- **THEN** 统一响应 `R.Code = 100005`，对应 HTTP Status 404

#### Scenario: 参数校验失败

- **WHEN** 请求参数未通过校验
- **THEN** 服务返回 `xerr.PARAM_ERROR`（Code=100002，HTTP 400）

### Requirement: 业务错误码按领域分段

业务错误码 SHALL 按领域分段维护，各段范围互不重叠：

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

#### Scenario: 新增领域错误码

- **WHEN** 开发者需要为某领域新增错误码
- **THEN** 必须在对应领域文件（如 `pkg/xerr/course.go`）中添加常量，且不越出该领域分段、不与其他分段重叠

### Requirement: 错误码使用方式

logic 层 SHALL 统一通过 `xerr` 包返回错误，错误 Msg 统一在 `xerr` 包内维护：

- 业务校验失败：`xerr.NewErrCode(xerr.USER_NOT_FOUND)` 形式返回 Code
- 包装底层错误：`xerr.NewErrCode(xerr.DB_ERROR).WithError(err)`
- 自定义提示：`xerr.NewErrMsg("自定义错误提示")`

#### Scenario: logic 返回错误码

- **WHEN** logic 中查询用户为 nil
- **THEN** 返回 `xerr.NewErrCode(xerr.USER_NOT_FOUND)`，禁止裸返回 error 或 `fmt.Errorf` 包装
- **AND** 错误文案由 `xerr` 包按 Code 统一给出，logic 不自行拼装 Msg

