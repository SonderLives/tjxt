> 版本：v1.0 | 更新：2026-08-05 | 来源：`apps/auth/rpc/internal/logic/*.go`

---

# Auth Business Rules

## 1. RBAC 角色管理

**核心规则**：角色分两种类型——固定角色（type=0，如 admin）和自定义角色（type=1）。

| 规则 | 说明 |
|------|------|
| 角色 code 唯一性 | 新增时校验 code 不重复，更新时可排除自身 |
| 固定角色不可改代号 | type=0 的角色修改 code 时拒绝 |
| 固定角色不可删除 | type=0 的角色 DeleteRole 直接返回冲突错误 |
| 已分配的固定角色不可删除 | 需先解除所有账户绑定 |
| 软删除 | 删除时 `deleted=1` 并失效缓存 |

```
流程（SaveRole）:
  1. 校验 code 非空、name 非空
  2. 校验 code 唯一性（ExistsByCode）
  3. 有 id → 更新；无 id → 新增（生成 ID）
  4. 更新时：type=0 且 code 变更则拒绝
```

## 2. 菜单树管理

**核心规则**：菜单通过 `parent_id` 自引用构建树形结构。

| 规则 | 说明 |
|------|------|
| 父菜单不能是自身 | `id == parentId` 时拒绝 |
| 子菜单存在不可删 | 删除菜单前检查 `CountChildren > 0` |
| 删除级联清理 | 删除菜单时同时删除其下所有权限、角色关联 |
| has_children 同步 | 增/删/移动菜单后调用 `SyncHasChildren` 更新父节点和自身 |
| 排序优先级 | 值越小越前，默认 127 |

```
流程（GetMenuTree）:
  1. 查所有菜单（FindAll）
  2. 构建 id→节点 map
  3. 按 parentId 拼接子节点
  4. 提取 root（parentId <= 0）
  5. 按 priority 排序
```

## 3. 权限管理

**核心规则**：权限依附于菜单，描述具体的 API 访问规则。

| 规则 | 说明 |
|------|------|
| method 统一大写 | `strings.ToUpper` 规范化 |
| method 和 uri 不可空 | 基础校验 |
| 删除权限时解绑所有角色 | `DeleteByPrivilegeId` 清理关联 |
| 软删除 | 逻辑删除 `deleted=1` |

## 4. 角色-菜单/权限分配（原子替换）

**核心规则**：分配操作使用 Replace 模式，即先清除旧关联再插入新关联，原子性保证。

| 规则 | 说明 |
|------|------|
| 角色存在性校验 | 分配前校验角色存在且未删除 |
| 资源存在性校验 | 批量分配菜单/权限时，逐 ID 校验是否存在 |
| Replace 语义 | `ReplaceByRoleId` 先 delete 再 batch insert |

```
流程（SaveRoleMenus）:
  1. 校验 roleId 有效
  2. 查角色，确认存在且未删除
  3. 若 menuIds 非空，逐 ID 校验存在性
  4. ReplaceByRoleId(roleId, menuIds) 原子替换
```

## 5. 账户-角色绑定

**核心规则**：账户可绑定多角色，绑定关系使用 Replace 语义。

| 规则 | 说明 |
|------|------|
| 账户 ID 必须 > 0 | 基础校验 |
| 每个角色 ID 必须 > 0 且存在 | 逐角色校验 |
| 已删除角色不可绑定 | role.deleted = 1 时拒绝 |
| Replace 替换 | 同一账户旧角色全部清除再插入新角色 |

## 6. JWT 令牌签发

**核心规则**：使用 `pkg/auth` 工具包的 `Sign` 方法签发 JWT。

| 配置项 | 默认值 | 说明 |
|--------|--------|------|
| `AccessExpire` | 7200 秒 (2 小时) | 访问令牌有效期 |
| `RefreshExpire` | 604800 秒 (7 天) | 刷新令牌有效期 |
| AccessSecret | 配置项 `AccessSecret` | JWT 签名密钥 |
| RefreshSecret | 配置项 `RefreshSecret` | 刷新令牌独立密钥 |

```
流程（SignToken）:
  1. expireSec <= 0 → 取 Jwt.AccessExpire → 再缺省 7200
  2. authutil.Sign(userId, roleCode, accessSecret, expire) → accessToken
  3. authutil.Sign(userId, roleCode, refreshSecret, refreshExpire) → refreshToken
  4. 返回 accessToken, refreshToken, expiresAt (unix秒)
```

## 7. 登录记录

| 规则 | 说明 |
|------|------|
| 容错插入 | 插入失败只 log 错误，不抛异常（`l.Errorf`） |
| 按 userId 分页 | 登录记录按用户 ID 分组查询 |

## 状态说明

### 角色类型

| 值 | 含义 | 限制 |
|----|------|------|
| 0 | 固定角色 | 不可改代号、不可删除 |
| 1 | 自定义角色 | 无限制 |

### 菜单 has_children

| 值 | 含义 |
|----|------|
| 0 | 无子菜单 |
| 1 | 有子菜单 |