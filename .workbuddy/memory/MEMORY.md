# 天机学堂 tjxt 项目记忆

## 架构
- go-zero v1.10.3 微服务，go.work 工作区管理 6 个服务模块 + pkg。
- 服务模块路径：`module tjxt/apps/<svc>`（每服务一个 go.mod，rpc 子目录无独立 go.mod）。
- 公共库 `tjxt/pkg`，各服务通过 `replace tjxt/pkg => ../../pkg` 引用。
- 命名风格 gozero（无下划线），goctl 默认即此风格。

## 服务清单与端口
| 服务 | API | RPC | DB |
|------|-----|-----|----|
| user | ✅8808 | ✅user.rpc:8082 | tj_user(user/user_detail) |
| course | ✅8812 | ✅course.rpc | tj_course |
| trade | ✅8810 | ✅trade.rpc | tj_trade |
| learning | ✅8888 | ✅learning.rpc:9000 | tj_learning |
| pay | ✅8811 | ✅pay.rpc:8081 | tj_pay |
| media | ✅8813 | ✅media.rpc:8087 | tj_media(file/media) |

## go-zero 开发规范（用户强约束）
- 优先 goctl 生成，只手写 logic/custom model/业务扩展。禁止手写骨架/handler/types/routes/pb。
- 开发顺序：DB表 → Model(goctl) → Proto → RPC生成 → .api → API生成 → Logic。
- Model 已存在则扩展 custom model，不重新生成覆盖。
- goctl .api 语法不支持 `any`，用 `interface{}` 代替。
- 统一响应：handler 用 `result.Write(w,r,data,err)`（pkg/response.R{Code,Msg,RequestId,Data any}），不用生成的 Result 类型。
- 统一分页：PageRequest{PageNo,PageSize}。
- JWT：@server jwt:Auth。

## MySQL（开发环境）
host 127.0.0.1 port 3306 user root pass 0000，库 tj_<domain>。

## 关键文件位置
- DDL：sql/ddl/tj_<domain>.sql（纯 DDL，供 goctl model 用）。
- 迁移：sql/migration/tj_<domain>.sql（含数据）。
- zero-skills：.cursor/skills/zero-skills/ 和 .opencode/skills/zero-skills/（go-zero AI 最佳实践）。
- Model 放 `apps/<svc>/rpc/internal/model/`（custom `*model.go` + 生成 `*model_gen.go` + vars.go）。
- RPC 客户端在 `apps/<svc>/rpc/<svc>/<svc>.go`（goctl 1.10 命名）。

## goctl 命令（本机已验证）
- API: `cd apps/X/api && goctl api go --api X.api --dir . --style gozero`
- RPC: `cd apps/X/rpc && goctl rpc protoc X.proto --go_out=. --go-grpc_out=. --zrpc_out=. --style gozero`
- Model: `goctl model mysql ddl -src sql/ddl/tj_X.sql -dir apps/X/rpc/internal/model -cache --style gozero`
- goctl api go 不覆盖已有 handler/logic/svc/config，仅刷新 types.go/routes.go。
- goctl rpc protoc 不覆盖已有 logic/svc/config，仅刷新 pb/server，client 默认生成(-c true)。

## 服务实现状态（logic 真实实现进度，2026-08-05）
| 服务 | API | RPC | 状态 |
|------|-----|-----|------|
| auth | 18/18 | 19/19 | ✅ |
| user | 13/13 | 15/15 | ✅ |
| pay | 16/16 | 17/17 | ✅ |
| promotion | 16/16 | 16/16 | ✅ |
| remark | 2/2 | 2/2 | ✅ |
| course | 38/38 | 38/38 | ✅（2026-08-05 完成） |
| trade | 0/36 | 0/37 | ⬜ |
| message | 0/18 | 0/19 | ⬜ |
| learning | 0/9 | 0/11 | ⬜ |
| media | 0/10 | 0/10 | ⬜ |
| exam | 0/7 | 0/7 | ⬜ |
| search | 0/2 | 0/2 | ⬜ |
| data | 0/6 | 0/6 | ⬜（目录结构异常待重建）|
合计 210/390（53.8%），6/13 服务完成。下一优先：media（course 依赖素材，须先补对象存储配置）。

## 批量化实现 stub 服务的可复用工作流（course 已验证）
对「goctl 已生成骨架、logic 全为 todo 占位」的服务，用此流程高效落地：
1. **先补 Model 层**：生成 model 只有 CRUD；业务需 List/Count/Page/FindByX/Upsert/级联 Delete，写 custom `*model.go` 扩展（不覆盖 `_gen.go`）。
2. **致命坑**：自定义 model 嵌入 `CachedConn`，**只能用** `m.QueryRowsNoCacheCtx`/`QueryRowNoCacheCtx`/`ExecNoCacheCtx`；`m.QueryRowsCtx` 被遮蔽会编译失败。
3. **建黄金标准**：挑一个完整域（如 course 的 category，8 RPC+8 API）手写实现，作为质量样板。
4. **写 IMPL_CONTRACT.md**：汇总「可调用方法清单 / 字段类型映射 / 共享 helper / 参考样板文件 / 已知缺口 / 硬性规则」，作为并行子代理的唯一事实源（任务结束即删）。
5. **并行子代理**：按域拆分（如 5 RPC + 4 API 共 9 个 general-purpose 子代理），每个只改非重叠文件集；**严禁子代理跑 go build**（并发写导致冲突），统一由主代理编译修复。
6. **主代理统一编译**：`cd apps/X/rpc && go build ./...` 与 `cd apps/X/api && go build ./...` 必须 rc=0。
7. **抽样核对 + 文档同步**：抽查 2–3 个关键 logic 确为真实实现（非 stub）；同步 repowiki business-rules.md 与 implementation-status.md。
8. 跨域字段（销量/评分/老师详情来自其他服务）当前服务无列则填 0/留空，并在文档「已知缺口」标注未接线的 RPC client。
