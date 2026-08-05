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
