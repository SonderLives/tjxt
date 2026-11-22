param(
    [string]$ProjectRoot = "."
)

# 与 Makefile build 目标保持一致：13 服务 x {api,rpc} = 26 个二进制
$services = @(
    @{path="api/auth"; name="auth-api"},
    @{path="rpc/auth"; name="auth-rpc"},
    @{path="api/course"; name="course-api"},
    @{path="rpc/course"; name="course-rpc"},
    @{path="api/data"; name="data-api"},
    @{path="rpc/data"; name="data-rpc"},
    @{path="api/exam"; name="exam-api"},
    @{path="rpc/exam"; name="exam-rpc"},
    @{path="api/learning"; name="learning-api"},
    @{path="rpc/learning"; name="learning-rpc"},
    @{path="api/media"; name="media-api"},
    @{path="rpc/media"; name="media-rpc"},
    @{path="api/message"; name="message-api"},
    @{path="rpc/message"; name="message-rpc"},
    @{path="api/pay"; name="pay-api"},
    @{path="rpc/pay"; name="pay-rpc"},
    @{path="api/promotion"; name="promotion-api"},
    @{path="rpc/promotion"; name="promotion-rpc"},
    @{path="api/remark"; name="remark-api"},
    @{path="rpc/remark"; name="remark-rpc"},
    @{path="api/search"; name="search-api"},
    @{path="rpc/search"; name="search-rpc"},
    @{path="api/trade"; name="trade-api"},
    @{path="rpc/trade"; name="trade-rpc"},
    @{path="api/user"; name="user-api"},
    @{path="rpc/user"; name="user-rpc"}
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
