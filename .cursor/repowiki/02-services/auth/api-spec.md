# auth 服务 HTTP API 接口规格

> 来源：`docs/tjxt.openapi.json` (聚合文档) | 最后同步：2026-08-05

## 接口清单

| 方法 | 路径 | 摘要 | 认证 | 请求体 | 响应体 | 权限标签 |
|------|------|------|------|--------|--------|----------|
| POST | /accounts/admin/login | 管理端登录并获取token | 否 | LoginFormDTO | R{data: R%C2%ABstring%C2%BB} |  |
| POST | /accounts/login | 登录并获取token | 否 | LoginFormDTO | R{data: R%C2%ABstring%C2%BB} |  |
| POST | /accounts/logout | 退出登录 | 否 | - | R{data: R} |  |
| GET | /accounts/refresh | 刷新token | 否 | - | R{data: R%C2%ABstring%C2%BB} |  |
| GET | /menus | 查询菜单，按照多级菜单组成树结构 | 否 | - | R{data: R%C2%ABList%C2%ABMenuOptionVO%C2%BB%C2%BB} |  |
| POST | /menus | 新增菜单 | 否 | MenuDTO | R{data: R} |  |
| GET | /menus/me | 查询我的菜单，按照多级菜单组成树结构 | 否 | - | R{data: R%C2%ABList%C2%ABMenuOptionVO%C2%BB%C2%BB} |  |
| GET | /menus/parent/{pid} | 根据父菜单id查询子菜单 | 否 | - | R{data: R%C2%ABList%C2%ABMenuOptionVO%C2%BB%C2%BB} |  |
| DELETE | /menus/role/{roleId} | 解除角色的菜单权限 | 否 | bindRoleMenusUsingPOSTMenuids | R{data: R} |  |
| POST | /menus/role/{roleId} | 绑定角色与菜单权限 | 否 | bindRoleMenusUsingPOSTMenuids | R{data: R} |  |
| DELETE | /menus/{id} | 根据id删除菜单 | 否 | - | R{data: R} |  |
| GET | /menus/{id} | 根据id查询菜单 | 否 | - | R{data: R%C2%ABMenuOptionVO%C2%BB} |  |
| PUT | /menus/{id} | 更新菜单 | 否 | MenuDTO | R{data: R} |  |
| GET | /privileges | 分页查询所有权限 | 否 | - | R{data: R%C2%ABPageDTO%C2%ABPrivilegeDTO%C2%BB%C2%BB} |  |
| POST | /privileges | 新增权限 | 否 | PrivilegeDTO | R{data: R%C2%ABPrivilegeDTO%C2%BB} |  |
| GET | /privileges/options/{menuId} | 查询菜单下的所有权限，作为下拉选框菜单 | 否 | - | R{data: R%C2%ABList%C2%ABPrivilegeOptionVO%C2%BB%C2%BB} |  |
| DELETE | /privileges/role/{roleId} | 解除角色的API权限 | 否 | bindRolePrivilegesUsingPOSTPrivilegeids | R{data: R} |  |
| POST | /privileges/role/{roleId} | 绑定角色与API权限 | 否 | bindRolePrivilegesUsingPOSTPrivilegeids | R{data: R} |  |
| GET | /privileges/roles/{roleId}/{menuId} | 查询菜单下的权限列表，某个角色的权限 | 否 | - | R{data: R%C2%ABList%C2%ABPrivilegeOptionVO%C2%BB%C2%BB} |  |
| DELETE | /privileges/{id} | 删除权限 | 否 | - | R{data: R} |  |
| PUT | /privileges/{id} | 修改权限 | 否 | PrivilegeDTO | R{data: R%C2%ABPrivilegeDTO%C2%BB} |  |
| GET | /roles | 查询员工角色列表 | 否 | - | R{data: R%C2%ABList%C2%ABRoleDTO%C2%BB%C2%BB} |  |
| POST | /roles | 新增角色 | 否 | RoleDTO | R{data: R%C2%ABRoleDTO%C2%BB} |  |
| GET | /roles/list | 查询员工角色列表 | 否 | - | R{data: R%C2%ABList%C2%ABRoleDTO%C2%BB%C2%BB} |  |
| DELETE | /roles/{id} | 删除角色信息 | 否 | - | R{data: R} |  |
| GET | /roles/{id} | 根据id查询角色 | 否 | - | R{data: R%C2%ABRoleDTO%C2%BB} |  |
| PUT | /roles/{id} | 修改角色信息 | 否 | RoleDTO | R{data: R} |  |

## 统一约定（引用全局规范）
- 响应格式：`pkg/response.R{Code,Msg,RequestId,Data any}`
- 分页：`PageRequest{PageNo,PageSize}` → `PageResponse{Total,List,PageNo,PageSize}`
- 错误码：见 [共享错误码表](../03-shared/error-codes.md)