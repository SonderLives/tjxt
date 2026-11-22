param(
    [string]$ProjectRoot = "."
)

# 与 Makefile build 目标保持一致：13 api + 5 保留 rpc + 2 mq consumer = 20 个二进制
# （8 个折叠服务 auth/exam/media/message/remark/search/trade/data 的 rpc 进程壳已删除，
#   业务实现在各自 app/<svc>/api/backend；trade/search 的 MQ 消费者独立为 app/mq/*）
$services = @(
    @{path="app/auth/api";          name="auth-api"},
    @{path="app/course/api";        name="course-api"},
    @{path="app/data/api";          name="data-api"},
    @{path="app/exam/api";          name="exam-api"},
    @{path="app/learning/api";      name="learning-api"},
    @{path="app/media/api";         name="media-api"},
    @{path="app/message/api";       name="message-api"},
    @{path="app/pay/api";           name="pay-api"},
    @{path="app/promotion/api";     name="promotion-api"},
    @{path="app/remark/api";        name="remark-api"},
    @{path="app/search/api";        name="search-api"},
    @{path="app/trade/api";         name="trade-api"},
    @{path="app/user/api";          name="user-api"},
    @{path="app/user/rpc";          name="user-rpc"},
    @{path="app/course/rpc";        name="course-rpc"},
    @{path="app/learning/rpc";      name="learning-rpc"},
    @{path="app/pay/rpc";           name="pay-rpc"},
    @{path="app/promotion/rpc";     name="promotion-rpc"},
    @{path="app/mq/trade-mq";   name="trade-mq"},
    @{path="app/mq/search-mq";  name="search-mq"}
)

$ProjectRoot = (Resolve-Path $ProjectRoot).Path
$binDir = Join-Path $ProjectRoot "bin"
if (-not (Test-Path $binDir)) {
    New-Item -ItemType Directory -Path $binDir | Out-Null
}

Set-Location $ProjectRoot
$fail = 0
foreach ($svc in $services) {
    $outPath = Join-Path $binDir "$($svc.name).exe"
    Write-Host "  构建 $($svc.path) -> bin/$($svc.name).exe..."
    go build -o $outPath "./$($svc.path)"
    if ($LASTEXITCODE -ne 0) { $fail = 1 }
}

if ($fail) {
    Write-Host "构建失败" -ForegroundColor Red
    exit 1
}
Write-Host "构建完成！20 个二进制文件在 bin/ 目录" -ForegroundColor Green
