# remark 服务 HTTP API 接口规格

> 来源：`apps/remark/api/remark.api` | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /likes | 点赞或取消点赞 | 是 | LikeRecordFormReq | OkVO | like |
| GET | /likes/list | 批量查询当前用户已点赞的业务id | 是 | LikeListReq (query) | LikeListResp | like |

## 请求/响应结构

### LikeRecordFormReq

| 字段 | 类型 | 位置 | 说明 |
|------|------|------|------|
| `bizId` | int64 | json | 被点赞的业务 id |
| `bizType` | string | json | 业务类型，如 reply / note / question |
| `liked` | bool | json | true=点赞，false=取消点赞 |

### OkVO

| 字段 | 类型 | 说明 |
|------|------|------|
| `success` | bool | 固定返回 true |

### LikeListReq

| 字段 | 类型 | 位置 | 说明 |
|------|------|------|------|
| `bizType` | string | form | 业务类型 |
| `bizIds` | []string | form | 业务 id 列表，服务端逐个 `strconv.ParseInt` 转换，非法或 `<= 0` 的直接丢弃 |

### LikeListResp

| 字段 | 类型 | 说明 |
|------|------|------|
| `likedBizIds` | []int64 | 入参 `bizIds` 中当前用户已点赞（liked=1）的子集 |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)
- 鉴权：`remark.api` 中 `@server(jwt: Auth, group: like)`，`routes.go` 由 `rest.WithJwt(serverCtx.Config.Auth.AccessSecret)` 统一包裹，两个接口均需携带有效 JWT；用户身份由 `auth.UserIdFromCtx(ctx)` 从上下文取出
