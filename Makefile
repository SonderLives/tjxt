.PHONY: help init api rpc model proto build test fmt lint clean docker-up docker-down docker-logs migrate-up migrate-down generate verify run-course run-trade run-learning run-pay run-course-rpc run-trade-rpc run-learning-rpc run-pay-rpc run-media-rpc run-all-rpc sync tidy

# Default target
.DEFAULT_GOAL := help

# Go related variables
GOCTL := goctl
GO := go

# Project services (relative to apps/)
# Only services with .api files
API_SERVICES := course/api learning/api pay/api trade/api
# Only services with .proto files
RPC_SERVICES := course/rpc learning/rpc media/rpc pay/rpc trade/rpc

# Help
help: ## 显示帮助信息
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@powershell -Command "Get-Content $(MAKEFILE_LIST) | Where-Object { $$_ -match '^[a-zA-Z_-]+:.*##' } | ForEach-Object { $$parts = $$_ -split ':.*##'; Write-Host \"  $$($$parts[0]).PadRight(20) $$parts[1]\" -ForegroundColor Cyan }"

# Initialize project - download dependencies and generate code
init: ## 初始化项目 - 整理模块并生成代码
	@echo "初始化项目..."
	$(GO) work sync
	@$(MAKE) generate
	@echo "项目初始化完成！"

# Generate API code from .api files
api: ## 从 .api 文件生成 API handler、logic、types 等
	@echo "生成 API 代码..."
	@powershell -Command "$(API_SERVICES) | ForEach-Object { $$api = Get-ChildItem -Path \"apps/$$_/*.api\" -ErrorAction SilentlyContinue | Select-Object -First 1; if ($$api) { Write-Host \"  生成 $$_...\"; & $(GOCTL) api go -api $$api.FullName -dir \"apps/$$_\" -style gozero } else { Write-Host \"  跳过 $$_ (未找到 .api 文件)\" } }"
	@echo "API 生成完成！"

# Generate RPC code from .proto files
rpc: ## 从 .proto 文件生成 RPC server、client、pb
	@echo "生成 RPC 代码..."
	@powershell -Command "$(RPC_SERVICES) | ForEach-Object { $$proto = Get-ChildItem -Path \"apps/$$_/*.proto\" -ErrorAction SilentlyContinue | Select-Object -First 1; if ($$proto) { Write-Host \"  生成 $$_...\"; & $(GOCTL) rpc protoc -proto $$proto.FullName -go_out \"apps/$$_\" -go-grpc_out \"apps/$$_\" -zrpc_out \"apps/$$_\" } else { Write-Host \"  跳过 $$_ (未找到 .proto 文件)\" } }"
	@echo "RPC 生成完成！"

# Generate Model code from MySQL DDL
model: ## 从 MySQL DDL 生成 model 代码 (需要 MySQL 运行)
	@echo "从 DDL 生成 Model 代码..."
	@echo "注意: 确保 MySQL 在 127.0.0.1:3306 运行，用户:root 密码:0000"
	@powershell -Command "Get-ChildItem -Path sql/ddl/*.sql | ForEach-Object { $$db_name = $_.BaseName -replace 'tj_', ''; Write-Host \"  处理 $$db_name...\"; & $(GOCTL) model mysql ddl -src $_.FullName -dir \"./model/$$db_name\" -c }"
	@echo "Model 生成完成！"

# Alternative: Generate model from datasource (requires DB)
model-datasource: ## 从活跃的 MySQL 数据源生成 model (需要 MySQL 运行)
	@echo "从数据源生成 Model..."
	@powershell -Command "'user','course','trade','learning','pay','auth','media','exam','message','search','promotion','remark' | ForEach-Object { Write-Host \"  生成 model for $$_...\"; & $(GOCTL) model mysql datasource -url \"root:0000@tcp(127.0.0.1:3306)/tj_$$_\" -dir \"./model/$$_\" -cache -style gozero }"
	@echo "数据源 Model 生成完成！"

# Generate proto from existing go files (if needed)
proto: ## 生成 proto 文件 (预留)
	@echo "Proto 生成暂未实现"

# Build all services
build: ## 构建所有服务
	@echo "构建所有服务..."
	@powershell -Command "if (-not (Test-Path bin)) { New-Item -ItemType Directory -Path bin | Out-Null }; $$services = @(@{path='apps/course/api'; name='course-api'}, @{path='apps/learning/api'; name='learning-api'}, @{path='apps/pay/api'; name='pay-api'}, @{path='apps/trade/api'; name='trade-api'}, @{path='apps/course/rpc'; name='course-rpc'}, @{path='apps/learning/rpc'; name='learning-rpc'}, @{path='apps/media/rpc'; name='media-rpc'}, @{path='apps/pay/rpc'; name='pay-rpc'}, @{path='apps/trade/rpc'; name='trade-rpc'}); $$services | ForEach-Object { Write-Host \"  构建 $$($_.path) -> bin/$$($_.name)...\"; go build -o \"bin/$$($_.name)\" \"./$$($_.path)\" }"
	@echo "构建完成！二进制文件在 bin/ 目录"

# Run tests
test: ## 运行所有测试
	@echo "运行测试..."
	$(GO) test ./... -v

# Format code
fmt: ## 格式化所有 Go 代码
	@echo "格式化代码..."
	$(GO) fmt ./...
	@echo "格式化完成！"

# Lint code
lint: ## 检查所有 Go 代码
	@echo "代码检查..."
	@which golangci-lint > /dev/null && golangci-lint run ./... || echo "golangci-lint 未安装，跳过..."

# Clean generated files
clean: ## 清理生成的文件和构建产物
	@echo "清理中..."
	@powershell -Command "if (Test-Path bin) { Remove-Item -Recurse -Force bin }"
	@powershell -Command "$(API_SERVICES) $(RPC_SERVICES) | ForEach-Object { $$exe = Get-ChildItem -Path \"apps/$$_/*.exe\" -ErrorAction SilentlyContinue; if ($$exe) { $$exe | Remove-Item -Force } }"
	@echo "清理完成！"

# Docker compose
docker-up: ## 启动 docker compose 服务
	docker-compose up -d

docker-down: ## 停止 docker compose 服务
	docker-compose down

docker-logs: ## 查看 docker compose 日志
	docker-compose logs -f

# Database migration
migrate-up: ## 执行数据库迁移 (预留)
	@echo "迁移功能暂未实现"

migrate-down: ## 回滚数据库迁移 (预留)
	@echo "迁移功能暂未实现"

# Generate all code
generate: api rpc model ## 生成所有代码
	@echo "所有代码生成完成！"

# Verify project structure
verify: ## 验证 go-zero 项目结构
	@echo "验证项目结构..."
	@powershell -Command "$$apiServices = @('course/api','learning/api','pay/api','trade/api'); $$rpcServices = @('course/rpc','learning/rpc','media/rpc','pay/rpc','trade/rpc'); $$errors = 0; $$apiServices | ForEach-Object { $$api = Get-ChildItem -Path \"apps/$$_/*.api\" -ErrorAction SilentlyContinue; if (-not $$api) { Write-Error \"错误: $$_ 缺少 .api 文件\"; $$errors++ } if (-not (Test-Path \"apps/$$_/internal\")) { Write-Error \"错误: $$_ 缺少 internal 目录\"; $$errors++ } }; $$rpcServices | ForEach-Object { $$proto = Get-ChildItem -Path \"apps/$$_/*.proto\" -ErrorAction SilentlyContinue; if (-not $$proto) { Write-Error \"错误: $$_ 缺少 .proto 文件\"; $$errors++ } if (-not (Test-Path \"apps/$$_/internal\")) { Write-Error \"错误: $$_ 缺少 internal 目录\"; $$errors++ } }; if ($$errors -gt 0) { exit 1 } else { Write-Host '项目结构验证通过!' -ForegroundColor Green }"

# Run individual services
run-course: ## 运行课程服务
	cd apps/course/api && $(GO) run . -f etc/course-api.yaml

run-trade: ## 运行交易服务
	cd apps/trade/api && $(GO) run . -f etc/trade-api.yaml

run-learning: ## 运行学习服务
	cd apps/learning/api && $(GO) run . -f etc/learning-api.yaml

run-pay: ## 运行支付服务
	cd apps/pay/api && $(GO) run . -f etc/pay-api.yaml

# Run RPC services
run-course-rpc: ## 运行课程 RPC
	cd apps/course/rpc && $(GO) run . -f etc/course-rpc.yaml

run-trade-rpc: ## 运行交易 RPC
	cd apps/trade/rpc && $(GO) run . -f etc/trade-rpc.yaml

run-learning-rpc: ## 运行学习 RPC
	cd apps/learning/rpc && $(GO) run . -f etc/learning-rpc.yaml

run-pay-rpc: ## 运行支付 RPC
	cd apps/pay/rpc && $(GO) run . -f etc/pay-rpc.yaml

run-media-rpc: ## 运行媒资 RPC
	cd apps/media/rpc && $(GO) run . -f etc/media-rpc.yaml

# Run all RPC services in background
run-all-rpc: run-course-rpc run-trade-rpc run-learning-rpc run-pay-rpc run-media-rpc ## 运行所有 RPC 服务

# Sync go.work dependencies
sync: ## 同步 go.work 依赖
	$(GO) work sync

# Tidy all modules
tidy: ## 整理所有模块依赖
	@powershell -Command "$$services = @('course/api','learning/api','pay/api','trade/api','course/rpc','learning/rpc','media/rpc','pay/rpc','trade/rpc'); $$services | ForEach-Object { $$mod = \"apps/$$_/go.mod\"; if (Test-Path $$mod) { Write-Host \"  整理 $$_...\"; Set-Location \"apps/$_\"; & $(GO) mod tidy; Set-Location ../.. } }"
	@powershell -Command "if (Test-Path go.mod) { & $(GO) mod tidy }"
	@echo "依赖整理完成！"