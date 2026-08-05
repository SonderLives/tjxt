# search 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| GET | /interests | 查询我的兴趣爱好 | 否 | - | R{data: R%C2%ABList%C2%ABCategoryBasicDTO%C2%BB%C2%BB} |  |
| POST | /interests | 新增兴趣爱好 | 否 | - | R{data: R} |  |
| GET | /interests/{id}/courses | 根据二级分类id查询课程TOP10 | 否 | - | R{data: R%C2%ABList%C2%ABCourseVO%C2%BB%C2%BB} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)