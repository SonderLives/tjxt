param(
    [string]$ProjectRoot = "."
)

# 与 Makefile build 目标保持一致：13 服务 x {api,rpc} = 26 个二进制
# data 是嵌套结构：api/data、rpc/data 才是模块根
$services = @(
    @{path="apps/auth/api"; name="auth-api"},
    @{path="apps/auth/rpc"; name="auth-rpc"},
    @{path="apps/course/api"; name="course-api"},
    @{path="apps/course/rpc"; name="course-rpc"},
    @{path="apps/data/api/data"; name="data-api"},
    @{path="apps/data/rpc/data"; name="data-rpc"},
    @{path="apps/exam/api"; name="exam-api"},
    @{path="apps/exam/rpc"; name="exam-rpc"},
    @{path="apps/learning/api"; name="learning-api"},
    @{path="apps/learning/rpc"; name="learning-rpc"},
    @{path="apps/media/api"; name="media-api"},
    @{path="apps/media/rpc"; name="media-rpc"},
    @{path="apps/message/api"; name="message-api"},
    @{path="apps/message/rpc"; name="message-rpc"},
    @{path="apps/pay/api"; name="pay-api"},
    @{path="apps/pay/rpc"; name="pay-rpc"},
    @{path="apps/promotion/api"; name="promotion-api"},
    @{path="apps/promotion/rpc"; name="promotion-rpc"},
    @{path="apps/remark/api"; name="remark-api"},
    @{path="apps/remark/rpc"; name="remark-rpc"},
    @{path="apps/search/api"; name="search-api"},
    @{path="apps/search/rpc"; name="search-rpc"},
    @{path="apps/trade/api"; name="trade-api"},
    @{path="apps/trade/rpc"; name="trade-rpc"},
    @{path="apps/user/api"; name="user-api"},
    @{path="apps/user/rpc"; name="user-rpc"}
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
Write-Host "构建完成！26 个二进制文件在 bin/ 目录" -ForegroundColor Green
