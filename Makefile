.DEFAULT_GOAL := help

GOCTL := goctl
GO    := go

# 所有 API 服务的启动目录（go run 在这个目录里执行）
# data 是嵌套结构：api/data 才是服务根
API_DIRS := apps/auth/api apps/course/api apps/data/api/data apps/exam/api apps/learning/api apps/media/api apps/message/api apps/pay/api apps/search/api apps/trade/api apps/user/api

# 所有 RPC 服务的启动目录
RPC_DIRS := apps/auth/rpc apps/course/rpc apps/data/rpc/data apps/exam/rpc apps/learning/rpc apps/media/rpc apps/message/rpc apps/pay/rpc apps/search/rpc apps/trade/rpc apps/user/rpc

# 服务名（用于二进制文件命名、run 目标解析）
SERVICES := auth course data exam learning media message pay search trade user

.PHONY: help init generate api rpc model build test fmt lint clean \
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
#   一般服务：apps/<svc>/api/<svc>.api -> apps/<svc>/api
#   data   ：apps/data/api/data.api   -> apps/data/api/data
api: ## 根据 .api 重新生成 handler/logic/types（覆盖生成产物，logic 手改注意备份）
	@powershell -NoProfile -Command "$$jobs = @( \
	  @{ api='apps/auth/api/auth.api';       dir='apps/auth/api' }, \
	  @{ api='apps/course/api/course.api';   dir='apps/course/api' }, \
	  @{ api='apps/data/api/data.api';       dir='apps/data/api/data' }, \
	  @{ api='apps/exam/api/exam.api';       dir='apps/exam/api' }, \
	  @{ api='apps/learning/api/learning.api'; dir='apps/learning/api' }, \
	  @{ api='apps/media/api/media.api';     dir='apps/media/api' }, \
	  @{ api='apps/message/api/message.api'; dir='apps/message/api' }, \
	  @{ api='apps/pay/api/pay.api';         dir='apps/pay/api' }, \
	  @{ api='apps/search/api/search.api';   dir='apps/search/api' }, \
	  @{ api='apps/trade/api/trade.api';     dir='apps/trade/api' }, \
	  @{ api='apps/user/api/user.api';       dir='apps/user/api' } ); \
	  $$jobs | ForEach-Object { Write-Host (\"  api  -> {0}\" -f $$_.dir); & $(GOCTL) api go -api $$_.api -dir $$_.dir -style gozero; if ($$LASTEXITCODE -ne 0) { exit 1 } }"

# goctl rpc protoc：在 .proto 所在目录下执行
rpc: ## 根据 .proto 重新生成 pb/server/client/logic
	@powershell -NoProfile -Command "$$jobs = @( \
	  'apps/auth/rpc/auth.proto', 'apps/course/rpc/course.proto', 'apps/data/rpc/data/data.proto', \
	  'apps/exam/rpc/exam.proto', 'apps/learning/rpc/learning.proto', 'apps/media/rpc/media.proto', \
	  'apps/message/rpc/message.proto', 'apps/pay/rpc/pay.proto', 'apps/search/rpc/search.proto', \
	  'apps/trade/rpc/trade.proto', 'apps/user/rpc/user.proto' ); \
	  $$jobs | ForEach-Object { $$d = Split-Path $$_ -Parent; $$f = Split-Path $$_ -Leaf; Write-Host (\"  rpc  -> {0}\" -f $$d); Push-Location $$d; & $(GOCTL) rpc protoc $$f --go_out=. --go-grpc_out=. --zrpc_out=. --client=true -m; $$code=$$LASTEXITCODE; Pop-Location; if ($$code -ne 0) { exit 1 } }"

# 从 DDL 生成 model（--cache 走 Redis）
model: ## 根据 sql/ddl/*.sql 生成 model（输出到 ./model/<db>/）
	@powershell -NoProfile -Command "Get-ChildItem sql/ddl/*.sql | ForEach-Object { $$db = $$_.BaseName -replace '^tj_',''; Write-Host (\"  model -> model/{0}  from {1}\" -f $$db, $$_.Name); & $(GOCTL) model mysql ddl -src $$_.FullName -dir (\"./model/{0}\" -f $$db) -c -style gozero; if ($$LASTEXITCODE -ne 0) { exit 1 } }"

generate: api rpc ## 一键重新生成 api + rpc

# ---------- 构建 / 测试 ----------

build: ## 构建全部 22 个可执行文件到 bin/
	@powershell -NoProfile -Command "New-Item -ItemType Directory -Force -Path bin | Out-Null; \
	  $$jobs = @( \
	    @{d='apps/auth/api';       n='auth-api'},       @{d='apps/auth/rpc';       n='auth-rpc'}, \
	    @{d='apps/course/api';     n='course-api'},     @{d='apps/course/rpc';     n='course-rpc'}, \
	    @{d='apps/data/api/data';  n='data-api'},       @{d='apps/data/rpc/data';  n='data-rpc'}, \
	    @{d='apps/exam/api';       n='exam-api'},       @{d='apps/exam/rpc';       n='exam-rpc'}, \
	    @{d='apps/learning/api';   n='learning-api'},   @{d='apps/learning/rpc';   n='learning-rpc'}, \
	    @{d='apps/media/api';      n='media-api'},      @{d='apps/media/rpc';      n='media-rpc'}, \
	    @{d='apps/message/api';    n='message-api'},    @{d='apps/message/rpc';    n='message-rpc'}, \
	    @{d='apps/pay/api';        n='pay-api'},        @{d='apps/pay/rpc';        n='pay-rpc'}, \
	    @{d='apps/search/api';     n='search-api'},     @{d='apps/search/rpc';     n='search-rpc'}, \
	    @{d='apps/trade/api';      n='trade-api'},      @{d='apps/trade/rpc';      n='trade-rpc'}, \
	    @{d='apps/user/api';       n='user-api'},       @{d='apps/user/rpc';       n='user-rpc'} ); \
	  $$fail = 0; \
	  $$jobs | ForEach-Object { Write-Host (\"  build {0,-22} -> bin/{1}.exe\" -f $$_.d, $$_.n); & $(GO) build -o (\"bin/{0}.exe\" -f $$_.n) (\"./{0}\" -f $$_.d); if ($$LASTEXITCODE -ne 0) { $$fail = 1 } }; \
	  if ($$fail) { exit 1 }; Write-Host 'done -> bin/' -ForegroundColor Green"

test: ## go test 所有模块
	@powershell -NoProfile -Command "$(API_DIRS) $(RPC_DIRS) | ForEach-Object { Push-Location $$_; Write-Host (\"  test {0}\" -f $$_); & $(GO) test ./...; Pop-Location }"

fmt: ## go fmt 所有模块
	@powershell -NoProfile -Command "$(API_DIRS) $(RPC_DIRS) | ForEach-Object { Push-Location $$_; & $(GO) fmt ./...; Pop-Location }"

lint: ## golangci-lint（如果已安装）
	@powershell -NoProfile -Command "if (Get-Command golangci-lint -ErrorAction SilentlyContinue) { golangci-lint run ./... } else { Write-Host 'golangci-lint 未安装，跳过' -ForegroundColor Yellow }"

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
	  @('apps/auth/api/auth.api','apps/course/api/course.api','apps/data/api/data.api','apps/exam/api/exam.api','apps/learning/api/learning.api','apps/media/api/media.api','apps/message/api/message.api','apps/pay/api/pay.api','apps/search/api/search.api','apps/trade/api/trade.api','apps/user/api/user.api') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host (\"missing api: {0}\" -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('apps/auth/rpc/auth.proto','apps/course/rpc/course.proto','apps/data/rpc/data/data.proto','apps/exam/rpc/exam.proto','apps/learning/rpc/learning.proto','apps/media/rpc/media.proto','apps/message/rpc/message.proto','apps/pay/rpc/pay.proto','apps/search/rpc/search.proto','apps/trade/rpc/trade.proto','apps/user/rpc/user.proto') | ForEach-Object { if (-not (Test-Path $$_)) { Write-Host (\"missing proto: {0}\" -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  @('$(API_DIRS)' -split ' ' + '$(RPC_DIRS)' -split ' ') | ForEach-Object { if (-not (Test-Path (\"{0}/internal\" -f $$_))) { Write-Host (\"missing internal: {0}\" -f $$_) -ForegroundColor Red; $$miss++ } }; \
	  if ($$miss -gt 0) { exit 1 } else { Write-Host 'structure ok' -ForegroundColor Green }"

sync: ## go work sync
	$(GO) work sync

tidy: ## 对所有 go.mod 执行 go mod tidy
	@powershell -NoProfile -Command "Push-Location pkg; & $(GO) mod tidy; Pop-Location; $(API_DIRS) $(RPC_DIRS) | ForEach-Object { Push-Location $$_; Write-Host (\"  tidy {0}\" -f $$_); & $(GO) mod tidy; Pop-Location }"

# ---------- 运行 ----------

# 模板：make run-<svc>        -> 启动 <svc>-api
#        make run-<svc>-rpc   -> 启动 <svc>-rpc
# data 目录是 apps/data/api/data 和 apps/data/rpc/data，yaml 名固定

run-auth: ; $(GO) run ./apps/auth/api -f apps/auth/api/etc/auth-api.yaml
run-course: ; $(GO) run ./apps/course/api -f apps/course/api/etc/course-api.yaml
run-data: ; $(GO) run ./apps/data/api/data -f apps/data/api/data/etc/data-api.yaml
run-exam: ; $(GO) run ./apps/exam/api -f apps/exam/api/etc/exam-api.yaml
run-learning: ; $(GO) run ./apps/learning/api -f apps/learning/api/etc/learning-api.yaml
run-media: ; $(GO) run ./apps/media/api -f apps/media/api/etc/media-api.yaml
run-message: ; $(GO) run ./apps/message/api -f apps/message/api/etc/message-api.yaml
run-pay: ; $(GO) run ./apps/pay/api -f apps/pay/api/etc/pay-api.yaml
run-search: ; $(GO) run ./apps/search/api -f apps/search/api/etc/search-api.yaml
run-trade: ; $(GO) run ./apps/trade/api -f apps/trade/api/etc/trade-api.yaml
run-user: ; $(GO) run ./apps/user/api -f apps/user/api/etc/user-api.yaml

run-auth-rpc: ; $(GO) run ./apps/auth/rpc -f apps/auth/rpc/etc/auth.yaml
run-course-rpc: ; $(GO) run ./apps/course/rpc -f apps/course/rpc/etc/course.yaml
run-data-rpc: ; $(GO) run ./apps/data/rpc/data -f apps/data/rpc/data/etc/data.yaml
run-exam-rpc: ; $(GO) run ./apps/exam/rpc -f apps/exam/rpc/etc/exam.yaml
run-learning-rpc: ; $(GO) run ./apps/learning/rpc -f apps/learning/rpc/etc/learning.yaml
run-media-rpc: ; $(GO) run ./apps/media/rpc -f apps/media/rpc/etc/media.yaml
run-message-rpc: ; $(GO) run ./apps/message/rpc -f apps/message/rpc/etc/message.yaml
run-pay-rpc: ; $(GO) run ./apps/pay/rpc -f apps/pay/rpc/etc/pay.yaml
run-search-rpc: ; $(GO) run ./apps/search/rpc -f apps/search/rpc/etc/search.yaml
run-trade-rpc: ; $(GO) run ./apps/trade/rpc -f apps/trade/rpc/etc/trade.yaml
run-user-rpc: ; $(GO) run ./apps/user/rpc -f apps/user/rpc/etc/user.yaml

# 后台并行拉起所有 RPC（每个进程留在自己的 PowerShell 窗口里）
run-all-rpc: ## 后台启动所有 RPC 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',(\"make run-{0}-rpc\" -f $$_) }"

# 后台并行拉起所有 API
run-all: ## 后台启动所有 API 服务（各自独立窗口）
	@powershell -NoProfile -Command "@('auth','course','data','exam','learning','media','message','pay','search','trade','user') | ForEach-Object { Start-Process powershell -ArgumentList '-NoExit','-Command',(\"make run-{0}\" -f $$_) }"