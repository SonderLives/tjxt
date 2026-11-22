.DEFAULT_GOAL := help

GOCTL := goctl
GO    := go

# goctl 自定义模板目录（api/rpc 生成统一走该模板，保证生成产物稳定）
GOCTL_TPL_HOME := tools/goctl-templates

# Windows 下强制用 cmd 作为 recipe shell。
# 否则 PATH 里有 sh.exe（如 Git Bash）时 make 会改用 sh，把 PS 片段里的 $$jobs、$$_
# 当 sh 变量展开，打碎所有 powershell 目标（build/api/rpc/test/lint/verify/tidy/fmt/clean）。
ifeq ($(OS),Windows_NT)
SHELL := cmd.exe
.SHELLFLAGS := /c
endif

# 所有 API 服务的启动目录（go run 在这个目录里执行）
API_DIRS := app/auth/api app/course/api app/data/api app/exam/api app/learning/api app/media/api app/message/api app/pay/api app/promotion/api app/remark/api app/search/api app/trade/api app/user/api

# 保留 RPC 服务的启动目录（折叠后仅剩 5 个有跨服务调用方的服务）
RPC_DIRS := app/user/rpc app/course/rpc app/learning/rpc app/pay/rpc app/promotion/rpc

# 折叠服务（RPC 进程壳已删除，业务实现在 app/<svc>/api/backend；pb 类型走 make pb / buf generate）
FOLDED_SVCS := auth exam media message remark search trade data

# MQ 消费者服务目录（自 rpc 进程独立出来的两个后台消费者）
MQ_DIRS := app/mq/trade-mq app/mq/search-mq

# MQ 消费者服务名（用于 run-*-mq 目标）
MQ_SERVICES := trade-mq search-mq

# 服务名（用于二进制文件命名、run 目标解析）
SERVICES := auth course data exam learning media message pay promotion remark search trade user

# ---------- 数据库迁移（pressly/goose） ----------
# 12 个业务库一个目录一个库：sql/migrations/<db>/（00001_schema.sql + 00002_seed.sql）。
# thirdparty/（gogs/nacos/seata/xxl_job）不归 goose 管理；业务库需先建库（sql/bootstrap.sql
# 或 docker-up 首次启动自动建库），再执行 make migrate-up。
MIGRATE_DBS := auth course learning exam media message pay promotion remark search trade user
MIGRATE_HOST ?= 127.0.0.1:3306
MIGRATE_USER ?= root
MIGRATE_PASS ?= $(if $(TJXT_DB_PASS),$(TJXT_DB_PASS),0000)
MIGRATE_DSN  := $(MIGRATE_USER):$(MIGRATE_PASS)@tcp($(MIGRATE_HOST))
# migrate-down 单次步数（goose down 本身无步数参数，recipe 内按单步循环 N 次）
MIGRATE_STEPS ?= 1

# 保留 RPC 服务名（用于 run-*-rpc / run-all-rpc）
RPC_SERVICES := user course learning pay promotion

.PHONY: help init generate api rpc pb model check-gen build test fmt lint ci clean \
        docker-up docker-down docker-logs verify sync tidy \
        migrate-up migrate-down migrate-status \
        buf-lint buf-breaking \
        run-all $(SERVICES:%=run-%) $(RPC_SERVICES:%=run-%-rpc) $(MQ_SERVICES:%=run-%)

help: ## 显示帮助
	@echo Usage: make [target]
	@echo.
	@powershell -NoProfile -Command "Get-Content $(MAKEFILE_LIST) | Where-Object { $$_ -match '^[a-zA-Z][a-zA-Z0-9_-]*:.*## ' } | ForEach-Object { $$p = $$_ -split ':.*## ',2; Write-Host ('  {0,-22} {1}' -f $$p[0], $$p[1]) -ForegroundColor Cyan }"

init: ## go work sync 并 tidy 所有模块
	$(GO) work sync
	@$(MAKE) tidy

# ---------- 代码生成 ----------

# goctl api go：在 .api 所在目录下生成，输出目录与 .api 同级
#   所有服务：app/<svc>/api/<svc>.api -> app/<svc>/api
api: ## 根据 .api 重新生成 handler/logic/types（覆盖生成产物，logic 手改注意备份）
	@powershell -NoProfile -Command "$$jobs = @( \
	  @{ api='app/auth/api/auth.api';       dir='app/auth/api' }, \
	  @{ api='app/course/api/course.api';   dir='app/course/api' }, \
	  @{ api='app/data/api/data.api';       dir='app/data/api' }, \
	  @{ api='app/exam/api/exam.api';       dir='app/exam/api' }, \
	  @{ api='app/learning/api/learning.api'; dir='app/learning/api' }, \
	  @{ api='app/media/api/media.api';     dir='app/media/api' }, \
	  @{ api='app/message/api/message.api'; dir='app/message/api' }, \
	  @{ api='app/pay/api/pay.api';         dir='app/pay/api' }, \
	  @{ api='app/promotion/api/promotion.api'; dir='app/promotion/api' }, \
	  @{ api='app/remark/api/remark.api';   dir='app/remark/api' }, \
	  @{ api='app/search/api/search.api';   dir='app/search/api' }, \
	  @{ api='app/trade/api/trade.api';     dir='app/trade/api' }, \
	  @{ api='app/user/api/user.api';       dir='app/user/api' } ); \
	  $$jobs | ForEach-Object { Write-Host ('  api  -> {0}' -f $$_.dir); & $(GOCTL) api go -api $$_.api -dir $$_.dir -style gozero -home $(GOCTL_TPL_HOME); if ($$LASTEXITCODE -ne 0) { exit 1 } }"

# goctl rpc protoc：在 .proto 所在目录下执行
# 仅 5 个保留 RPC 服务；8 个折叠服务的 pb 类型生成走 pb 目标（buf generate）。
rpc: ## 根据 .proto 重新生成 pb/server/client/logic（5 个保留 RPC 服务）
	@powershell -NoProfile -Command "$$jobs = @( \
	  'app/course/rpc/course.proto', 'app/learning/rpc/learning.proto', \
	  'app/pay/rpc/pay.proto', 'app/promotion/rpc/promotion.proto', 'app/user/rpc/user.proto' ); \
	  $$jobs | ForEach-Object { $$d = Split-Path $$_ -Parent; $$f = Split-Path $$_ -Leaf; Write-Host ('  rpc  -> {0}' -f $$d); Push-Location $$d; & $(GOCTL) rpc protoc $$f --go_out=. --go-grpc_out=. --zrpc_out=. --client=true -m --home $(GOCTL_TPL_HOME); $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# buf generate：为 8 个折叠服务再生成 pb 类型（app/<svc>/api/backend/pb）。
# 它们的 rpc 进程壳已删、goctl rpc protoc 无法只生成类型，故走 buf；
# buf 向上找到仓库根 buf.yaml workspace，只对当前模块生成，输出由各 proto 的
# go_package（./backend/pb）锚定到 backend/pb/。模板见仓库根 buf.gen.yaml（仅 protoc-gen-go）。
pb: ## buf generate 折叠服务的 pb 类型（app/<svc>/api/backend/pb）
	@powershell -NoProfile -Command "@('$(FOLDED_SVCS)') -split ' ' | ForEach-Object { Write-Host ('  pb   -> app/{0}/api' -f $$_); Push-Location ('app/{0}/api' -f $$_); & buf generate --template ../../../buf.gen.yaml; $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# 从迁移 schema 生成 model（--cache 走 Redis）。
# 生成源为 sql/migrations/<db>/00001_schema.sql（与真实库表集合一致）；三张表 goctl 1.10.1
# 无法解析——exam/trade 的 undo_log（无主键，Seata AT 基础表，无对应 model）与 message 的
# notice_task_target（复合主键）——生成前把对应 CREATE TABLE 块剥除到临时文件再生成
# （等价原 ddl/tj_exam_business.sql、tj_message_model.sql 变体文件的作用，两者已随 sql/ddl/ 退役）。
model: ## 根据 sql/migrations/<db>/00001_schema.sql 生成 model（输出到 ./model/<db>/；exam/trade/message 自动剥除 goctl 不解析的表）
	@powershell -NoProfile -Command "$$exclude = @{ exam = 'undo_log'; trade = 'undo_log'; message = 'notice_task_target' }; \
	  Get-ChildItem sql/migrations/*/00001_schema.sql | ForEach-Object { \
	    $$db = Split-Path (Split-Path $$_.FullName -Parent) -Leaf; $$src = $$_.FullName; \
	    if ($$exclude.ContainsKey($$db)) { \
	      $$tbl = $$exclude[$$db]; $$out = New-Object 'System.Collections.Generic.List[string]'; $$skip = $$false; \
	      foreach ($$l in [System.IO.File]::ReadAllLines($$src)) { \
	        if ($$skip) { if ($$l -match '^\)') { $$skip = $$false } else { continue } } \
	        elseif ($$l -match ('^\s*CREATE\s+TABLE\s+[\x60]?' + $$tbl + '[\x60]?\s*\(')) { $$skip = $$true } \
	        else { $$out.Add($$l) } }; \
	      $$src = Join-Path $$env:TEMP ('goose_model_' + $$db + '.sql'); \
	      [System.IO.File]::WriteAllLines($$src, $$out) }; \
	    Write-Host ('  model -> model/{0}  from {1}' -f $$db, $$src); \
	    & $(GOCTL) model mysql ddl -src $$src -dir ('./model/{0}' -f $$db) -c -style gozero; if ($$LASTEXITCODE -ne 0) { exit 1 } }"

generate: api rpc pb ## 一键重新生成 api + rpc + pb

check-gen: ## 重生成 API 与 pb 类型并校验与仓库零漂移（api types 与 pb 类型双守卫）
	$(MAKE) api
	$(MAKE) pb
	git diff --exit-code -- app

# ---------- 构建 / 测试 ----------

build: ## 构建全部 20 个可执行文件到 bin/（13 api + 5 rpc + 2 mq consumer）
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path bin | Out-Null; \
	  $$jobs = @( \
	    @{d='app/auth/api';       n='auth-api'},       @{d='app/course/api';     n='course-api'}, \
	    @{d='app/data/api';       n='data-api'},       @{d='app/exam/api';       n='exam-api'}, \
	    @{d='app/learning/api';   n='learning-api'},   @{d='app/media/api';      n='media-api'}, \
	    @{d='app/message/api';    n='message-api'},    @{d='app/pay/api';        n='pay-api'}, \
	    @{d='app/promotion/api';  n='promotion-api'},  @{d='app/remark/api';     n='remark-api'}, \
	    @{d='app/search/api';     n='search-api'},     @{d='app/trade/api';      n='trade-api'}, \
	    @{d='app/user/api';       n='user-api'}, \
	    @{d='app/user/rpc';       n='user-rpc'},       @{d='app/course/rpc';     n='course-rpc'}, \
	    @{d='app/learning/rpc';   n='learning-rpc'},   @{d='app/pay/rpc';        n='pay-rpc'}, \
	    @{d='app/promotion/rpc';  n='promotion-rpc'}, \
	    @{d='app/mq/trade-mq';  n='trade-mq'}, \
	    @{d='app/mq/search-mq'; n='search-mq'} ); \
	  $$fail = 0; \
	  $$jobs | ForEach-Object { Write-Host ('  build {0,-22} -> bin/{1}.exe' -f $$_.d, $$_.n); & $(GO) build -o ('bin/{0}.exe' -f $$_.n) ('./{0}' -f $$_.d); if ($$LASTEXITCODE -ne 0) { $$fail = 1 } }; \
	  if ($$fail) { exit 1 }; Write-Host 'done -> bin/' -ForegroundColor Green"

test: ## go test 所有模块（含 pkg 与 mq consumer）
	@powershell -NoProfile -Command "('pkg $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  test {0}' -f $$_); & $(GO) test ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

fmt: ## go fmt 所有模块
	@powershell -NoProfile -Command "('pkg $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; & $(GO) fmt ./...; Pop-Location }"

lint: ## golangci-lint 逐模块检查（未安装或发现问题即失败退出）
	@powershell -NoProfile -Command "if (-not (Get-Command golangci-lint -ErrorAction SilentlyContinue)) { Write-Host 'golangci-lint 未安装：go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest 或 scoop install golangci-lint' -ForegroundColor Red; exit 1 }; ('pkg $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  lint {0}' -f $$_); & golangci-lint run ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

ci: verify lint test build ## 一条龙门禁：结构校验 + lint + 测试 + 全量构建（CI 与本地同一入口）

test-integration: ## 集成测试（依赖本地 MySQL/Redis/RabbitMQ，需先 docker-up；不可达时自动跳过）
	@powershell -NoProfile -Command "('pkg app/trade/api app/pay/rpc app/learning/rpc') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  itest {0}' -f $$_); & $(GO) test -tags integration -count=1 ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

clean: ## 删除 bin/ 与各服务下遗留的 exe
	@powershell -NoProfile -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }; $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS) | ForEach-Object { Get-ChildItem -Path $$_ -Filter *.exe -ErrorAction SilentlyContinue | Remove-Item -Force }"

# ---------- 本地依赖 ----------

docker-up: ## 启动 MySQL/Redis/RabbitMQ/etcd
	docker compose up -d

docker-down: ## 停止并移除容器
	docker compose down

docker-logs: ## 跟踪依赖容器日志
	docker compose logs -f

# ---------- 数据库迁移（goose） ----------
# 循环 12 个业务库执行 goose；第三方基建库（gogs/nacos/seata/xxl_job）不在其列，
# 由 sql/thirdparty/ 脚本在 docker 首次启动时初始化。业务库不存在时先执行
# sql/bootstrap.sql（docker 首次启动已自动执行）再 migrate-up。
# 口令默认 0000，可用环境变量 TJXT_DB_PASS 覆盖：make migrate-up TJXT_DB_PASS=xxx。

migrate-up: ## goose up 全部 12 个业务库（sql/migrations/<db>/：schema + seed）
	@powershell -NoProfile -Command "if (-not (Get-Command goose -ErrorAction SilentlyContinue)) { Write-Host 'goose 未安装：go install github.com/pressly/goose/v3/cmd/goose@latest' -ForegroundColor Red; exit 1 }; \
	  $$dsn = '$(MIGRATE_DSN)'; \
	  @('$(MIGRATE_DBS)') -split ' ' | ForEach-Object { Write-Host ('  goose up  -> tj_{0}' -f $$_); & goose -dir ('sql/migrations/{0}' -f $$_) mysql ($$dsn + '/tj_' + $$_) up; if ($$LASTEXITCODE -ne 0) { exit 1 } }; \
	  Write-Host 'migrate up ok' -ForegroundColor Green"

migrate-down: ## goose down 每库回退 MIGRATE_STEPS 步（默认 1，goose down 即单步语义）
	@powershell -NoProfile -Command "if (-not (Get-Command goose -ErrorAction SilentlyContinue)) { Write-Host 'goose 未安装：go install github.com/pressly/goose/v3/cmd/goose@latest' -ForegroundColor Red; exit 1 }; \
	  $$dsn = '$(MIGRATE_DSN)'; $$s = '$(MIGRATE_STEPS)'.Trim(); if (-not $$s) { $$s = '1' }; $$steps = [Math]::Max(1, [int]$$s); \
	  @('$(MIGRATE_DBS)') -split ' ' | ForEach-Object { \
	    $$db = $$_; Write-Host ('  goose down -> tj_{0} x{1}' -f $$db, $$steps); \
	    1..$$steps | ForEach-Object { & goose -dir ('sql/migrations/{0}' -f $$db) mysql ($$dsn + '/tj_' + $$db) down; if ($$LASTEXITCODE -ne 0) { exit 1 } } }; \
	  Write-Host 'migrate down ok' -ForegroundColor Green"

migrate-status: ## goose status 全部 12 个业务库
	@powershell -NoProfile -Command "if (-not (Get-Command goose -ErrorAction SilentlyContinue)) { Write-Host 'goose 未安装：go install github.com/pressly/goose/v3/cmd/goose@latest' -ForegroundColor Red; exit 1 }; \
	  $$dsn = '$(MIGRATE_DSN)'; \
	  @('$(MIGRATE_DBS)') -split ' ' | ForEach-Object { Write-Host ('  goose status -> tj_{0}' -f $$_); & goose -dir ('sql/migrations/{0}' -f $$_) mysql ($$dsn + '/tj_' + $$_) status; if ($$LASTEXITCODE -ne 0) { exit 1 } }; \
	  Write-Host 'migrate status ok' -ForegroundColor Green"

# ---------- 校验 ----------

verify: ## 校验所有服务的 .api/.proto 与 internal/ 目录是否齐全
	@powershell -NoProfile -Command "$$miss = 0; \
	  @('app/auth/api/auth.api','app/course/api/course.api','app/data/api/data.api','app/exam/api/exam.api','app/learning/api/learning.api','app/media/api/media.api','app/message/api/message.api','app/pay/api/pay.api','app/promotion/api/promotion.api','app/remark/api/remark.api','app/search/api/search.api','app/trade/api/trade.api','app/user/api/user.api') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing api: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('app/auth/api/auth.proto','app/exam/api/exam.proto','app/media/api/media.proto','app/message/api/message.proto','app/remark/api/remark.proto','app/search/api/search.proto','app/trade/api/trade.proto','app/data/api/data.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing proto: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('app/user/rpc/user.proto','app/course/rpc/course.proto','app/learning/rpc/learning.proto','app/pay/rpc/pay.proto','app/promotion/rpc/promotion.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing proto: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('$(API_DIRS)' -split ' ' + '$(RPC_DIRS)' -split ' ') | ForEach-Object { if (-not (Test-Path ('{0}/internal' -f $$_))) { Write-Host ('missing internal: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  if ($$miss -gt 0) { exit 1 } else { Write-Host 'structure ok' -ForegroundColor Green }"

# ---------- proto 治理（buf） ----------
# buf 承担 lint/breaking 门禁与折叠服务 pb 类型生成（make pb）；
# 5 个保留服务的 pb/zrpc 生成仍走 make rpc（goctl）。
# 配置见仓库根 buf.yaml（v2 workspace：5 个 app/<svc>/rpc + 8 个折叠服务 app/<svc>/api 模块）。

buf-lint: ## buf lint 检查 13 个 proto 模块（规则与豁免见 buf.yaml）
	buf lint

# breaking 对比 master（本 worktree 场景直接读本地 .git）。
# against 侧配置双保险：master 仍是 service/<svc>/{api,rpc} 布局（上一轮迁移的产物）时，
# 用本分支指向 app/<svc>/... 的 buf.yaml 解释 master 树会报「模块下无 .proto 文件」，
# 此时回退 buf.master.yaml（按 master 现行 service/ 布局声明同一批模块）；
# 合入后 master 树即 app/ 布局，第一条命令直接通过，回退不会执行。
buf-breaking: ## buf breaking 破坏性变更门禁（对比 master）
	buf breaking --against ".git#branch=master" --against-config buf.yaml || buf breaking --against ".git#branch=master" --against-config buf.master.yaml

sync: ## go work sync
	$(GO) work sync

tidy: ## 对所有 go.mod 执行 go mod tidy
	@powershell -NoProfile -Command "Push-Location pkg; & $(GO) mod tidy; Pop-Location; $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS) | ForEach-Object { Push-Location $$_; Write-Host ('  tidy {0}' -f $$_); & $(GO) mod tidy; Pop-Location }"

# ---------- 运行 ----------

# 模板：make run-<svc>        -> 启动 <svc>-api
#        make run-<svc>-rpc   -> 启动 <svc>-rpc（仅 5 个保留 RPC 服务）
#        make run-<mq>        -> 启动 MQ 消费者（trade-mq / search-mq）

run-auth: ; $(GO) run ./app/auth/api -f app/auth/api/etc/auth-api.yaml
run-course: ; $(GO) run ./app/course/api -f app/course/api/etc/course-api.yaml
run-data: ; $(GO) run ./app/data/api -f app/data/api/etc/data-api.yaml
run-exam: ; $(GO) run ./app/exam/api -f app/exam/api/etc/exam-api.yaml
run-learning: ; $(GO) run ./app/learning/api -f app/learning/api/etc/learning-api.yaml
run-media: ; $(GO) run ./app/media/api -f app/media/api/etc/media-api.yaml
run-message: ; $(GO) run ./app/message/api -f app/message/api/etc/message-api.yaml
run-pay: ; $(GO) run ./app/pay/api -f app/pay/api/etc/pay-api.yaml
run-promotion: ; $(GO) run ./app/promotion/api -f app/promotion/api/etc/promotion-api.yaml
run-remark: ; $(GO) run ./app/remark/api -f app/remark/api/etc/remark-api.yaml
run-search: ; $(GO) run ./app/search/api -f app/search/api/etc/search-api.yaml
run-trade: ; $(GO) run ./app/trade/api -f app/trade/api/etc/trade-api.yaml
run-user: ; $(GO) run ./app/user/api -f app/user/api/etc/user-api.yaml

run-user-rpc: ; $(GO) run ./app/user/rpc -f app/user/rpc/etc/user.yaml
run-course-rpc: ; $(GO) run ./app/course/rpc -f app/course/rpc/etc/course.yaml
run-learning-rpc: ; $(GO) run ./app/learning/rpc -f app/learning/rpc/etc/learning.yaml
run-pay-rpc: ; $(GO) run ./app/pay/rpc -f app/pay/rpc/etc/pay.yaml
run-promotion-rpc: ; $(GO) run ./app/promotion/rpc -f app/promotion/rpc/etc/promotion.yaml

run-trade-mq: ; $(GO) run ./app/mq/trade-mq -f app/mq/trade-mq/etc/trade-mq.yaml
run-search-mq: ; $(GO) run ./app/mq/search-mq -f app/mq/search-mq/etc/search-mq.yaml

# 后台并行拉起全部后台进程：5 个保留 RPC + 2 个 MQ 消费者（每个进程留在自己的 PowerShell 窗口里）
run-all-rpc: ## 后台启动所有 RPC 与 MQ 消费者（各自独立窗口）
	@powershell -NoProfile -Command "@('user','course','learning','pay','promotion') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}-rpc' -f $$_) }; @('trade-mq','search-mq') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}' -f $$_) }"

# 后台并行拉起所有 API
run-all: ## 后台启动所有 API 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','promotion','remark','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}' -f $$_) }"
