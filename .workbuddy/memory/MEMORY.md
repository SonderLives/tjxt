# 天机学堂 tjxt 项目记忆

## 架构
- go-zero v1.10.3 微服务，go.work 工作区管理 13 个服务模块 + pkg（apps/{auth,course,data,exam,learning,media,message,pay,promotion,remark,search,trade,user}）。
- 服务模块路径：`module tjxt/apps/<svc>`（每服务一个 go.mod，rpc 子目录无独立 go.mod）。
- 公共库 `tjxt/pkg`，各服务通过 `replace tjxt/pkg => ../../pkg` 引用。
- 命名风格 gozero（无下划线），goctl 默认即此风格。

## 服务清单与端口（2026-08-06 统一重排为连续编号，无冲突）
API 端口 8801–8813 连续唯一；RPC 端口 8081–8093 连续唯一。改端口不影响调用方（RPC 走 etcd Key 发现，如 `course.rpc`）。

| 服务 | API | RPC(ListenOn) | etcd Key | DB |
|------|-----|-----|-----|----|
| user | 8801 | 8081 | user.rpc | tj_user(user/user_detail) |
| auth | 8802 | 8082 | auth.rpc | tj_user |
| course | 8803 | 8083 | course.rpc | tj_course |
| learning | 8804 | 8084 | learning.rpc | tj_learning |
| exam | 8805 | 8085 | exam.rpc | tj_learning |
| media | 8806 | 8086 | media.rpc | tj_media(file/media) |
| message | 8807 | 8087 | message.rpc | tj_message |
| pay | 8808 | 8088 | pay.rpc | tj_pay |
| trade | 8809 | 8089 | trade.rpc | tj_trade |
| search | 8810 | 8090 | search.rpc | tj_search |
| data | 8811 | 8091 | data.rpc | tj_data |
| promotion | 8812 | 8092 | promotion.rpc | tj_promotion |
| remark | 8813 | 8093 | remark.rpc | tj_remark |

- data 为嵌套模块，配置在 `apps/data/api/data/etc/data-api.yaml`、`apps/data/rpc/data/etc/data.yaml`（不是 `apps/data/api/etc/`）。
- 端口在 26 个 `apps/**/etc/*.yaml` + README.md + `.cursor/repowiki/**/*.md` + `docs/**/*.md` 中均已同步；批量改端口时务必用「旧→新映射的单趟 re.sub」，因新旧端口集合重叠，逐条回扫会连锁错改。

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

## 服务实现状态（2026-08-06 实测：go build 全模块通过 + 逻辑文件清点）
- 13/13 服务 API+RPC 契约 100% 落地（API 193 + RPC 201 = 394 接口，394 个 logic 文件 1:1 对应，全代码库 0 处 TODO/panic 占位）。
- 跨服务 RPC 依赖（已接线、编译通过）：trade→{course,pay}、search→course、learning→course。
- 事件驱动：trade 经 RabbitMQ 发布领域事件，Producer 优雅降级（MQ 不可用时为 nil 不阻塞启动）。
- 代码注释标明的「应有但未接线」集成：course→user(教师)、course→learning(课时)、trade→promotion(优惠券)、pay→真实支付网关。
- 待优化（有意桩，非未实现）：media 对象存储为 mock（COS/OSS 未接入）；pay 支付回调 URL 为 demo 占位；trade 优惠券未接入；course 教师/课时未接线；learning section_type 忽略。
- 结构异常：apps/data/api/go.mod 为孤儿模块（真实模块为 apps/data/api/data，已纳入 go.work），建议清理。

| 服务 | API | RPC | 编译 | 备注 |
|------|-----|-----|:----:|------|
| auth | 18/18 | 19/19 | ✅ | 已实现 |
| user | 13/13 | 15/15 | ✅ | 已实现 |
| course | 38/38 | 38/38 | ✅ | 已实现；教师(user)/课时(learning)未接线 |
| media | 10/10 | 10/10 | ✅ | 已实现；对象存储为 mock |
| learning | 9/9 | 11/11 | ✅ | 已实现；section_type 占位 |
| exam | 7/7 | 7/7 | ✅ | 已实现 |
| search | 4/4 | 4/4 | ✅ | 已实现；依赖 course |
| trade | 36/36 | 37/37 | ✅ | 已实现；依赖 course,pay；优惠券未接入 |
| pay | 16/16 | 17/17 | ✅ | 已实现；支付网关占位 URL |
| promotion | 16/16 | 16/16 | ✅ | 已实现 |
| message | 18/18 | 19/19 | ✅ | 已实现 |
| remark | 2/2 | 2/2 | ✅ | 已实现 |
| data | 6/6 | 6/6 | ✅* | 已实现（*真实模块 apps/data/api/data 编译通过；apps/data/api 为孤儿 go.mod）|
合计 394/394（100%）逻辑文件落地，13/13 真实模块编译通过。整体代码实现完成度 ≈100%，功能完备度 ≈97%（少量外部集成桩）。完整报告见 docs/completion-analysis.md。

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
