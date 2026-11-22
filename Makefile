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
API_DIRS := service/auth/api service/course/api service/data/api service/exam/api service/learning/api service/media/api service/message/api service/pay/api service/promotion/api service/remark/api service/search/api service/trade/api service/user/api

# 保留 RPC 服务的启动目录（折叠后仅剩 5 个有跨服务调用方的服务）
RPC_DIRS := service/user/rpc service/course/rpc service/learning/rpc service/pay/rpc service/promotion/rpc

# 折叠服务（RPC 进程壳已删除，业务实现在 service/<svc>/api/backend；pb 类型走 make pb / buf generate）
FOLDED_SVCS := auth exam media message remark search trade data

# MQ 消费者服务目录（自 rpc 进程独立出来的两个后台消费者）
MQ_DIRS := app/mq/trade-mq app/mq/search-mq

# MQ 消费者服务名（用于 run-*-mq 目标）
MQ_SERVICES := trade-mq search-mq

# 服务名（用于二进制文件命名、run 目标解析）
SERVICES := auth course data exam learning media message pay promotion remark search trade user

# 保留 RPC 服务名（用于 run-*-rpc / run-all-rpc）
RPC_SERVICES := user course learning pay promotion

.PHONY: help init generate api rpc pb model check-gen build test fmt lint ci clean \
        docker-up docker-down docker-logs verify sync tidy \
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
#   所有服务：service/<svc>/api/<svc>.api -> service/<svc>/api
api: ## 根据 .api 重新生成 handler/logic/types（覆盖生成产物，logic 手改注意备份）
	@powershell -NoProfile -Command "$$jobs = @( \
	  @{ api='service/auth/api/auth.api';       dir='service/auth/api' }, \
	  @{ api='service/course/api/course.api';   dir='service/course/api' }, \
	  @{ api='service/data/api/data.api';       dir='service/data/api' }, \
	  @{ api='service/exam/api/exam.api';       dir='service/exam/api' }, \
	  @{ api='service/learning/api/learning.api'; dir='service/learning/api' }, \
	  @{ api='service/media/api/media.api';     dir='service/media/api' }, \
	  @{ api='service/message/api/message.api'; dir='service/message/api' }, \
	  @{ api='service/pay/api/pay.api';         dir='service/pay/api' }, \
	  @{ api='service/promotion/api/promotion.api'; dir='service/promotion/api' }, \
	  @{ api='service/remark/api/remark.api';   dir='service/remark/api' }, \
	  @{ api='service/search/api/search.api';   dir='service/search/api' }, \
	  @{ api='service/trade/api/trade.api';     dir='service/trade/api' }, \
	  @{ api='service/user/api/user.api';       dir='service/user/api' } ); \
	  $$jobs | ForEach-Object { Write-Host ('  api  -> {0}' -f $$_.dir); & $(GOCTL) api go -api $$_.api -dir $$_.dir -style gozero -home $(GOCTL_TPL_HOME); if ($$LASTEXITCODE -ne 0) { exit 1 } }"

# goctl rpc protoc：在 .proto 所在目录下执行
# 仅 5 个保留 RPC 服务；8 个折叠服务的 pb 类型生成走 pb 目标（buf generate）。
rpc: ## 根据 .proto 重新生成 pb/server/client/logic（5 个保留 RPC 服务）
	@powershell -NoProfile -Command "$$jobs = @( \
	  'service/course/rpc/course.proto', 'service/learning/rpc/learning.proto', \
	  'service/pay/rpc/pay.proto', 'service/promotion/rpc/promotion.proto', 'service/user/rpc/user.proto' ); \
	  $$jobs | ForEach-Object { $$d = Split-Path $$_ -Parent; $$f = Split-Path $$_ -Leaf; Write-Host ('  rpc  -> {0}' -f $$d); Push-Location $$d; & $(GOCTL) rpc protoc $$f --go_out=. --go-grpc_out=. --zrpc_out=. --client=true -m --home $(GOCTL_TPL_HOME); $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# buf generate：为 8 个折叠服务再生成 pb 类型（service/<svc>/api/backend/pb）。
# 它们的 rpc 进程壳已删、goctl rpc protoc 无法只生成类型，故走 buf；
# buf 向上找到仓库根 buf.yaml workspace，只对当前模块生成，输出由各 proto 的
# go_package（./backend/pb）锚定到 backend/pb/。模板见仓库根 buf.gen.yaml（仅 protoc-gen-go）。
pb: ## buf generate 折叠服务的 pb 类型（service/<svc>/api/backend/pb）
	@powershell -NoProfile -Command "@('$(FOLDED_SVCS)') -split ' ' | ForEach-Object { Write-Host ('  pb   -> service/{0}/api' -f $$_); Push-Location ('service/{0}/api' -f $$_); & buf generate --template ../../../buf.gen.yaml; $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# 从 DDL 生成 model（--cache 走 Redis）
model: ## 根据 sql/ddl/*.sql 生成 model（输出到 ./model/<db>/）
	@powershell -NoProfile -Command "Get-ChildItem sql/ddl/*.sql | ForEach-Object { $$db = $$_.BaseName -replace '^tj_',''; Write-Host ('  model -> model/{0}  from {1}' -f $$db, $$_.Name); & $(GOCTL) model mysql ddl -src $$_.FullName -dir ('./model/{0}' -f $$db) -c -style gozero; if ($$LASTEXITCODE -ne 0) { exit 1 } }"

generate: api rpc pb ## 一键重新生成 api + rpc + pb

check-gen: ## 重生成 API 与 pb 类型并校验与仓库零漂移（api types 与 pb 类型双守卫）
	$(MAKE) api
	$(MAKE) pb
	git diff --exit-code -- service

# ---------- 构建 / 测试 ----------

build: ## 构建全部 20 个可执行文件到 bin/（13 api + 5 rpc + 2 mq consumer）
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path bin | Out-Null; \
	  $$jobs = @( \
	    @{d='service/auth/api';       n='auth-api'},       @{d='service/course/api';     n='course-api'}, \
	    @{d='service/data/api';       n='data-api'},       @{d='service/exam/api';       n='exam-api'}, \
	    @{d='service/learning/api';   n='learning-api'},   @{d='service/media/api';      n='media-api'}, \
	    @{d='service/message/api';    n='message-api'},    @{d='service/pay/api';        n='pay-api'}, \
	    @{d='service/promotion/api';  n='promotion-api'},  @{d='service/remark/api';     n='remark-api'}, \
	    @{d='service/search/api';     n='search-api'},     @{d='service/trade/api';      n='trade-api'}, \
	    @{d='service/user/api';       n='user-api'}, \
	    @{d='service/user/rpc';       n='user-rpc'},       @{d='service/course/rpc';     n='course-rpc'}, \
	    @{d='service/learning/rpc';   n='learning-rpc'},   @{d='service/pay/rpc';        n='pay-rpc'}, \
	    @{d='service/promotion/rpc';  n='promotion-rpc'}, \
	    @{d='app/mq/trade-mq';  n='trade-mq'}, \
	    @{d='app/mq/search-mq'; n='search-mq'} ); \
	  $$fail = 0; \
	  $$jobs | ForEach-Object { Write-Host ('  build {0,-22} -> bin/{1}.exe' -f $$_.d, $$_.n); & $(GO) build -o ('bin/{0}.exe' -f $$_.n) ('./{0}' -f $$_.d); if ($$LASTEXITCODE -ne 0) { $$fail = 1 } }; \
	  if ($$fail) { exit 1 }; Write-Host 'done -> bin/' -ForegroundColor Green"

test: ## go test 所有模块（含 common 与 mq consumer）
	@powershell -NoProfile -Command "('common $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  test {0}' -f $$_); & $(GO) test ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

fmt: ## go fmt 所有模块
	@powershell -NoProfile -Command "('common $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; & $(GO) fmt ./...; Pop-Location }"

lint: ## golangci-lint 逐模块检查（未安装或发现问题即失败退出）
	@powershell -NoProfile -Command "if (-not (Get-Command golangci-lint -ErrorAction SilentlyContinue)) { Write-Host 'golangci-lint 未安装：go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest 或 scoop install golangci-lint' -ForegroundColor Red; exit 1 }; ('common $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  lint {0}' -f $$_); & golangci-lint run ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

ci: verify lint test build ## 一条龙门禁：结构校验 + lint + 测试 + 全量构建（CI 与本地同一入口）

test-integration: ## 集成测试（依赖本地 MySQL/Redis/RabbitMQ，需先 docker-up；不可达时自动跳过）
	@powershell -NoProfile -Command "('common service/trade/api service/pay/rpc service/learning/rpc') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  itest {0}' -f $$_); & $(GO) test -tags integration -count=1 ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

clean: ## 删除 bin/ 与各服务下遗留的 exe
	@powershell -NoProfile -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }; $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS) | ForEach-Object { Get-ChildItem -Path $$_ -Filter *.exe -ErrorAction SilentlyContinue | Remove-Item -Force }"

# ---------- 本地依赖 ----------

docker-up: ## 启动 MySQL/Redis/RabbitMQ/etcd
	docker compose up -d

docker-down: ## 停止并移除容器
	docker compose down

docker-logs: ## 跟踪依赖容器日志
	docker compose logs -f

# ---------- 校验 ----------

verify: ## 校验所有服务的 .api/.proto 与 internal/ 目录是否齐全
	@powershell -NoProfile -Command "$$miss = 0; \
	  @('service/auth/api/auth.api','service/course/api/course.api','service/data/api/data.api','service/exam/api/exam.api','service/learning/api/learning.api','service/media/api/media.api','service/message/api/message.api','service/pay/api/pay.api','service/promotion/api/promotion.api','service/remark/api/remark.api','service/search/api/search.api','service/trade/api/trade.api','service/user/api/user.api') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing api: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('service/auth/api/auth.proto','service/exam/api/exam.proto','service/media/api/media.proto','service/message/api/message.proto','service/remark/api/remark.proto','service/search/api/search.proto','service/trade/api/trade.proto','service/data/api/data.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing proto: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('service/user/rpc/user.proto','service/course/rpc/course.proto','service/learning/rpc/learning.proto','service/pay/rpc/pay.proto','service/promotion/rpc/promotion.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing proto: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('$(API_DIRS)' -split ' ' + '$(RPC_DIRS)' -split ' ') | ForEach-Object { if (-not (Test-Path ('{0}/internal' -f $$_))) { Write-Host ('missing internal: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  if ($$miss -gt 0) { exit 1 } else { Write-Host 'structure ok' -ForegroundColor Green }"

# ---------- proto 治理（buf） ----------
# buf 承担 lint/breaking 门禁与折叠服务 pb 类型生成（make pb）；
# 5 个保留服务的 pb/zrpc 生成仍走 make rpc（goctl）。
# 配置见仓库根 buf.yaml（v2 workspace：5 个 service/<svc>/rpc + 8 个折叠服务 service/<svc>/api 模块）。

buf-lint: ## buf lint 检查 13 个 proto 模块（规则与豁免见 buf.yaml）
	buf lint

# breaking 对比 master（本 worktree 场景直接读本地 .git）；
# master 尚无 buf.yaml，需用 --against-config 复用同一 workspace 配置，
# 否则 against 侧会被当成单一模块（13 vs 1 报错）。
buf-breaking: ## buf breaking 破坏性变更门禁（对比 master）
	buf breaking --against ".git#branch=master" --against-config buf.yaml

sync: ## go work sync
	$(GO) work sync

tidy: ## 对所有 go.mod 执行 go mod tidy
	@powershell -NoProfile -Command "Push-Location common; & $(GO) mod tidy; Pop-Location; $(API_DIRS) $(RPC_DIRS) $(MQ_DIRS) | ForEach-Object { Push-Location $$_; Write-Host ('  tidy {0}' -f $$_); & $(GO) mod tidy; Pop-Location }"

# ---------- 运行 ----------

# 模板：make run-<svc>        -> 启动 <svc>-api
#        make run-<svc>-rpc   -> 启动 <svc>-rpc（仅 5 个保留 RPC 服务）
#        make run-<mq>        -> 启动 MQ 消费者（trade-mq / search-mq）

run-auth: ; $(GO) run ./service/auth/api -f service/auth/api/etc/auth-api.yaml
run-course: ; $(GO) run ./service/course/api -f service/course/api/etc/course-api.yaml
run-data: ; $(GO) run ./service/data/api -f service/data/api/etc/data-api.yaml
run-exam: ; $(GO) run ./service/exam/api -f service/exam/api/etc/exam-api.yaml
run-learning: ; $(GO) run ./service/learning/api -f service/learning/api/etc/learning-api.yaml
run-media: ; $(GO) run ./service/media/api -f service/media/api/etc/media-api.yaml
run-message: ; $(GO) run ./service/message/api -f service/message/api/etc/message-api.yaml
run-pay: ; $(GO) run ./service/pay/api -f service/pay/api/etc/pay-api.yaml
run-promotion: ; $(GO) run ./service/promotion/api -f service/promotion/api/etc/promotion-api.yaml
run-remark: ; $(GO) run ./service/remark/api -f service/remark/api/etc/remark-api.yaml
run-search: ; $(GO) run ./service/search/api -f service/search/api/etc/search-api.yaml
run-trade: ; $(GO) run ./service/trade/api -f service/trade/api/etc/trade-api.yaml
run-user: ; $(GO) run ./service/user/api -f service/user/api/etc/user-api.yaml

run-user-rpc: ; $(GO) run ./service/user/rpc -f service/user/rpc/etc/user.yaml
run-course-rpc: ; $(GO) run ./service/course/rpc -f service/course/rpc/etc/course.yaml
run-learning-rpc: ; $(GO) run ./service/learning/rpc -f service/learning/rpc/etc/learning.yaml
run-pay-rpc: ; $(GO) run ./service/pay/rpc -f service/pay/rpc/etc/pay.yaml
run-promotion-rpc: ; $(GO) run ./service/promotion/rpc -f service/promotion/rpc/etc/promotion.yaml

run-trade-mq: ; $(GO) run ./app/mq/trade-mq -f app/mq/trade-mq/etc/trade-mq.yaml
run-search-mq: ; $(GO) run ./app/mq/search-mq -f app/mq/search-mq/etc/search-mq.yaml

# 后台并行拉起全部后台进程：5 个保留 RPC + 2 个 MQ 消费者（每个进程留在自己的 PowerShell 窗口里）
run-all-rpc: ## 后台启动所有 RPC 与 MQ 消费者（各自独立窗口）
	@powershell -NoProfile -Command "@('user','course','learning','pay','promotion') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}-rpc' -f $$_) }; @('trade-mq','search-mq') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}' -f $$_) }"

# 后台并行拉起所有 API
run-all: ## 后台启动所有 API 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','promotion','remark','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}' -f $$_) }"
