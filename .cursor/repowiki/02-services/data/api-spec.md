# data 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /data/board | 看板数据获取 | 否 | - | R{data: R%C2%ABEchartsVO%C2%BB} |  |
| PUT | /data/board/set | 看板数据设置 | 否 | BoardDataSetDTO | R{data: R} |  |
| GET | /data/today | 获取今日数据 | 否 | - | R{data: R%C2%ABTodayDataVO%C2%BB} |  |
| PUT | /data/today/set | 设置线上数据 | 否 | TodayDataDTO | R{data: R} |  |
| GET | /data/top10 | top10数据获取 | 否 | - | R{data: R%C2%ABTop10DataVO%C2%BB} |  |
| PUT | /data/top10/set | 设置top10数据 | 否 | Top10DataSetDTO | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)