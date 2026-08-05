# user 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /users | 新增用户，一般是员工或教师 | 否 | UserDTO | R{data: R%C2%ABlong%C2%BB} |  |
| PUT | /users | 更新当前登录用户信息，可修改密码 | 否 | UserFormDTO | R{data: R} |  |
| GET | /users/checkCellphone | 检查用户手机号是否存在 | 否 | - | R{data: R%C2%ABboolean%C2%BB} |  |
| GET | /users/me | 获取当前登录用户信息 | 否 | - | R{data: R%C2%ABUserDetailVO%C2%BB} |  |
| GET | /users/{id} | 根据id查询用户信息 | 否 | - | R{data: R%C2%ABUserDTO%C2%BB} |  |
| PUT | /users/{id} | 更新用户信息 | 否 | UserDTO | R{data: R} |  |
| PUT | /users/{id}/password/default | 重置密码 | 否 | - | R{data: R} |  |
| PUT | /users/{id}/status/{status} | 修改用户状态, status=0为禁用，status=1为正常 | 否 | - | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)