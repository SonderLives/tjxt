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
API_DIRS := api/auth api/course api/data api/exam api/learning api/media api/message api/pay api/promotion api/remark api/search api/trade api/user

# 所有 RPC 服务的启动目录
RPC_DIRS := rpc/auth rpc/course rpc/data rpc/exam rpc/learning rpc/media rpc/message rpc/pay rpc/promotion rpc/remark rpc/search rpc/trade rpc/user

# 服务名（用于二进制文件命名、run 目标解析）
SERVICES := auth course data exam learning media message pay promotion remark search trade user

.PHONY: help init generate api rpc model check-gen build test fmt lint ci clean \
        docker-up docker-down docker-logs verify sync tidy \
        run-all $(SERVICES:%=run-%) $(SERVICES:%=run-%-rpc)

help: ## 显示帮助
	@echo Usage: make [target]
	@echo.
	@powershell -NoProfile -Command "Get-Content $(MAKEFILE_LIST) | Where-Object { $$_ -match '^[a-zA-Z][a-zA-Z0-9_-]*:.*## ' } | ForEach-Object { $$p = $$_ -split ':.*## ',2; Write-Host ('  {0,-22} {1}' -f $$p[0], $$p[1]) -ForegroundColor Cyan }"

init: ## go work sync 并 tidy 所有模块
	$(GO) work sync
	@$(MAKE) tidy

# ---------- 代码生成 ----------

# goctl api go：在 .api 所在目录下生成，输出目录与 .api 同级
#   所有服务：api/<svc>/<svc>.api -> api/<svc>
api: ## 根据 .api 重新生成 handler/logic/types（覆盖生成产物，logic 手改注意备份）
	@powershell -NoProfile -Command "$$jobs = @( \
	  @{ api='api/auth/auth.api';       dir='api/auth' }, \
	  @{ api='api/course/course.api';   dir='api/course' }, \
	  @{ api='api/data/data.api';       dir='api/data' }, \
	  @{ api='api/exam/exam.api';       dir='api/exam' }, \
	  @{ api='api/learning/learning.api'; dir='api/learning' }, \
	  @{ api='api/media/media.api';     dir='api/media' }, \
	  @{ api='api/message/message.api'; dir='api/message' }, \
	  @{ api='api/pay/pay.api';         dir='api/pay' }, \
	  @{ api='api/promotion/promotion.api'; dir='api/promotion' }, \
	  @{ api='api/remark/remark.api';   dir='api/remark' }, \
	  @{ api='api/search/search.api';   dir='api/search' }, \
	  @{ api='api/trade/trade.api';     dir='api/trade' }, \
	  @{ api='api/user/user.api';       dir='api/user' } ); \
	  $$jobs | ForEach-Object { Write-Host ('  api  -> {0}' -f $$_.dir); & $(GOCTL) api go -api $$_.api -dir $$_.dir -style gozero -home $(GOCTL_TPL_HOME); if ($$LASTEXITCODE -ne 0) { exit 1 } }"

# goctl rpc protoc：在 .proto 所在目录下执行
rpc: ## 根据 .proto 重新生成 pb/server/client/logic
	@powershell -NoProfile -Command "$$jobs = @( \
	  'rpc/auth/auth.proto', 'rpc/course/course.proto', 'rpc/data/data.proto', \
	  'rpc/exam/exam.proto', 'rpc/learning/learning.proto', 'rpc/media/media.proto', \
	  'rpc/message/message.proto', 'rpc/pay/pay.proto', 'rpc/promotion/promotion.proto', \
	  'rpc/remark/remark.proto', 'rpc/search/search.proto', \
	  'rpc/trade/trade.proto', 'rpc/user/user.proto' ); \
	  $$jobs | ForEach-Object { $$d = Split-Path $$_ -Parent; $$f = Split-Path $$_ -Leaf; Write-Host ('  rpc  -> {0}' -f $$d); Push-Location $$d; & $(GOCTL) rpc protoc $$f --go_out=. --go-grpc_out=. --zrpc_out=. --client=true -m --home $(GOCTL_TPL_HOME); $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# 从 DDL 生成 model（--cache 走 Redis）
model: ## 根据 sql/ddl/*.sql 生成 model（输出到 ./model/<db>/）
	@powershell -NoProfile -Command "Get-ChildItem sql/ddl/*.sql | ForEach-Object { $$db = $$_.BaseName -replace '^tj_',''; Write-Host ('  model -> model/{0}  from {1}' -f $$db, $$_.Name); & $(GOCTL) model mysql ddl -src $$_.FullName -dir ('./model/{0}' -f $$db) -c -style gozero; if ($$LASTEXITCODE -ne 0) { exit 1 } }"

generate: api rpc ## 一键重新生成 api + rpc

check-gen: ## 重生成 API 并校验与仓库零漂移（契约一致性守卫）
	$(MAKE) api
	git diff --exit-code -- api

# ---------- 构建 / 测试 ----------

build: ## 构建全部 26 个可执行文件到 bin/
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path bin | Out-Null; \
	  $$jobs = @( \
	    @{d='api/auth';       n='auth-api'},       @{d='rpc/auth';       n='auth-rpc'}, \
	    @{d='api/course';     n='course-api'},     @{d='rpc/course';     n='course-rpc'}, \
	    @{d='api/data';       n='data-api'},       @{d='rpc/data';       n='data-rpc'}, \
	    @{d='api/exam';       n='exam-api'},       @{d='rpc/exam';       n='exam-rpc'}, \
	    @{d='api/learning';   n='learning-api'},   @{d='rpc/learning';   n='learning-rpc'}, \
	    @{d='api/media';      n='media-api'},      @{d='rpc/media';      n='media-rpc'}, \
	    @{d='api/message';    n='message-api'},    @{d='rpc/message';    n='message-rpc'}, \
	    @{d='api/pay';        n='pay-api'},        @{d='rpc/pay';        n='pay-rpc'}, \
	    @{d='api/promotion';  n='promotion-api'},  @{d='rpc/promotion';  n='promotion-rpc'}, \
	    @{d='api/remark';     n='remark-api'},     @{d='rpc/remark';     n='remark-rpc'}, \
	    @{d='api/search';     n='search-api'},     @{d='rpc/search';     n='search-rpc'}, \
	    @{d='api/trade';      n='trade-api'},      @{d='rpc/trade';      n='trade-rpc'}, \
	    @{d='api/user';       n='user-api'},       @{d='rpc/user';       n='user-rpc'} ); \
	  $$fail = 0; \
	  $$jobs | ForEach-Object { Write-Host ('  build {0,-22} -> bin/{1}.exe' -f $$_.d, $$_.n); & $(GO) build -o ('bin/{0}.exe' -f $$_.n) ('./{0}' -f $$_.d); if ($$LASTEXITCODE -ne 0) { $$fail = 1 } }; \
	  if ($$fail) { exit 1 }; Write-Host 'done -> bin/' -ForegroundColor Green"

test: ## go test 所有模块（含 pkg）
	@powershell -NoProfile -Command "('pkg $(API_DIRS) $(RPC_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  test {0}' -f $$_); & $(GO) test ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

fmt: ## go fmt 所有模块
	@powershell -NoProfile -Command "('pkg $(API_DIRS) $(RPC_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; & $(GO) fmt ./...; Pop-Location }"

lint: ## golangci-lint 逐模块检查（未安装或发现问题即失败退出）
	@powershell -NoProfile -Command "if (-not (Get-Command golangci-lint -ErrorAction SilentlyContinue)) { Write-Host 'golangci-lint 未安装：go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest 或 scoop install golangci-lint' -ForegroundColor Red; exit 1 }; ('pkg $(API_DIRS) $(RPC_DIRS)') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  lint {0}' -f $$_); & golangci-lint run ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

ci: verify lint test build ## 一条龙门禁：结构校验 + lint + 测试 + 全量构建（CI 与本地同一入口）

test-integration: ## 集成测试（依赖本地 MySQL/Redis/RabbitMQ，需先 docker-up；不可达时自动跳过）
	@powershell -NoProfile -Command "('pkg rpc/pay rpc/trade rpc/learning') -split ' ' | ForEach-Object { Push-Location $$_; Write-Host ('  itest {0}' -f $$_); & $(GO) test -tags integration -count=1 ./...; if ($$LASTEXITCODE -ne 0) { Pop-Location; exit 1 }; Pop-Location }"

clean: ## 删除 bin/ 与各服务下遗留的 exe
	@powershell -NoProfile -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }; $(API_DIRS) $(RPC_DIRS) | ForEach-Object { Get-ChildItem -Path $$_ -Filter *.exe -ErrorAction SilentlyContinue | Remove-Item -Force }"

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
	  @('api/auth/auth.api','api/course/course.api','api/data/data.api','api/exam/exam.api','api/learning/learning.api','api/media/media.api','api/message/message.api','api/pay/pay.api','api/promotion/promotion.api','api/remark/remark.api','api/search/search.api','api/trade/trade.api','api/user/user.api') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing api: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('rpc/auth/auth.proto','rpc/course/course.proto','rpc/data/data.proto','rpc/exam/exam.proto','rpc/learning/learning.proto','rpc/media/media.proto','rpc/message/message.proto','rpc/pay/pay.proto','rpc/promotion/promotion.proto','rpc/remark/remark.proto','rpc/search/search.proto','rpc/trade/trade.proto','rpc/user/user.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host ('missing proto: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('$(API_DIRS)' -split ' ' + '$(RPC_DIRS)' -split ' ') | ForEach-Object { if (-not (Test-Path ('{0}/internal' -f $$_))) { Write-Host ('missing internal: {0}' -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  if ($$miss -gt 0) { exit 1 } else { Write-Host 'structure ok' -ForegroundColor Green }"

sync: ## go work sync
	$(GO) work sync

tidy: ## 对所有 go.mod 执行 go mod tidy
	@powershell -NoProfile -Command "Push-Location pkg; & $(GO) mod tidy; Pop-Location; $(API_DIRS) $(RPC_DIRS) | ForEach-Object { Push-Location $$_; Write-Host ('  tidy {0}' -f $$_); & $(GO) mod tidy; Pop-Location }"

# ---------- 运行 ----------

# 模板：make run-<svc>        -> 启动 <svc>-api
#        make run-<svc>-rpc   -> 启动 <svc>-rpc

run-auth: ; $(GO) run ./api/auth -f api/auth/etc/auth-api.yaml
run-course: ; $(GO) run ./api/course -f api/course/etc/course-api.yaml
run-data: ; $(GO) run ./api/data -f api/data/etc/data-api.yaml
run-exam: ; $(GO) run ./api/exam -f api/exam/etc/exam-api.yaml
run-learning: ; $(GO) run ./api/learning -f api/learning/etc/learning-api.yaml
run-media: ; $(GO) run ./api/media -f api/media/etc/media-api.yaml
run-message: ; $(GO) run ./api/message -f api/message/etc/message-api.yaml
run-pay: ; $(GO) run ./api/pay -f api/pay/etc/pay-api.yaml
run-promotion: ; $(GO) run ./api/promotion -f api/promotion/etc/promotion-api.yaml
run-remark: ; $(GO) run ./api/remark -f api/remark/etc/remark-api.yaml
run-search: ; $(GO) run ./api/search -f api/search/etc/search-api.yaml
run-trade: ; $(GO) run ./api/trade -f api/trade/etc/trade-api.yaml
run-user: ; $(GO) run ./api/user -f api/user/etc/user-api.yaml

run-auth-rpc: ; $(GO) run ./rpc/auth -f rpc/auth/etc/auth.yaml
run-course-rpc: ; $(GO) run ./rpc/course -f rpc/course/etc/course.yaml
run-data-rpc: ; $(GO) run ./rpc/data -f rpc/data/etc/data.yaml
run-exam-rpc: ; $(GO) run ./rpc/exam -f rpc/exam/etc/exam.yaml
run-learning-rpc: ; $(GO) run ./rpc/learning -f rpc/learning/etc/learning.yaml
run-media-rpc: ; $(GO) run ./rpc/media -f rpc/media/etc/media.yaml
run-message-rpc: ; $(GO) run ./rpc/message -f rpc/message/etc/message.yaml
run-pay-rpc: ; $(GO) run ./rpc/pay -f rpc/pay/etc/pay.yaml
run-promotion-rpc: ; $(GO) run ./rpc/promotion -f rpc/promotion/etc/promotion.yaml
run-remark-rpc: ; $(GO) run ./rpc/remark -f rpc/remark/etc/remark.yaml
run-search-rpc: ; $(GO) run ./rpc/search -f rpc/search/etc/search.yaml
run-trade-rpc: ; $(GO) run ./rpc/trade -f rpc/trade/etc/trade.yaml
run-user-rpc: ; $(GO) run ./rpc/user -f rpc/user/etc/user.yaml

# 后台并行拉起所有 RPC（每个进程留在自己的 PowerShell 窗口里）
run-all-rpc: ## 后台启动所有 RPC 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','promotion','remark','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}-rpc' -f $$_) }"

# 后台并行拉起所有 API
run-all: ## 后台启动所有 API 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','promotion','remark','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',('make run-{0}' -f $$_) }"
