param(
    [string]$ProjectRoot = "."
)

# 与 Makefile build 目标保持一致：13 api + 5 保留 rpc + 2 mq consumer = 20 个二进制
# （8 个折叠服务 auth/exam/media/message/remark/search/trade/data 的 rpc 进程壳已删除，
#   业务实现在各自 api/<svc>/backend；trade/search 的 MQ 消费者独立为 app/mq/*）
$services = @(
    @{path="api/auth";          name="auth-api"},
    @{path="api/course";        name="course-api"},
    @{path="api/data";          name="data-api"},
    @{path="api/exam";          name="exam-api"},
    @{path="api/learning";      name="learning-api"},
    @{path="api/media";         name="media-api"},
    @{path="api/message";       name="message-api"},
    @{path="api/pay";           name="pay-api"},
    @{path="api/promotion";     name="promotion-api"},
    @{path="api/remark";        name="remark-api"},
    @{path="api/search";        name="search-api"},
    @{path="api/trade";         name="trade-api"},
    @{path="api/user";          name="user-api"},
    @{path="rpc/user";          name="user-rpc"},
    @{path="rpc/course";        name="course-rpc"},
    @{path="rpc/learning";      name="learning-rpc"},
    @{path="rpc/pay";           name="pay-rpc"},
    @{path="rpc/promotion";     name="promotion-rpc"},
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
