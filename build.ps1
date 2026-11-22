param(
    [string]$ProjectRoot = "."
)

# 与 Makefile build 目标保持一致：13 api + 5 保留 rpc + 2 mq consumer = 20 个二进制
# （8 个折叠服务 auth/exam/media/message/remark/search/trade/data 的 rpc 进程壳已删除，
#   业务实现在各自 service/<svc>/api/backend；trade/search 的 MQ 消费者独立为 app/mq/*）
$services = @(
    @{path="service/auth/api";          name="auth-api"},
    @{path="service/course/api";        name="course-api"},
    @{path="service/data/api";          name="data-api"},
    @{path="service/exam/api";          name="exam-api"},
    @{path="service/learning/api";      name="learning-api"},
    @{path="service/media/api";         name="media-api"},
    @{path="service/message/api";       name="message-api"},
    @{path="service/pay/api";           name="pay-api"},
    @{path="service/promotion/api";     name="promotion-api"},
    @{path="service/remark/api";        name="remark-api"},
    @{path="service/search/api";        name="search-api"},
    @{path="service/trade/api";         name="trade-api"},
    @{path="service/user/api";          name="user-api"},
    @{path="service/user/rpc";          name="user-rpc"},
    @{path="service/course/rpc";        name="course-rpc"},
    @{path="service/learning/rpc";      name="learning-rpc"},
    @{path="service/pay/rpc";           name="pay-rpc"},
    @{path="service/promotion/rpc";     name="promotion-rpc"},
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
