param(
    [string]$ProjectRoot = "."
)

$services = @(
    @{path="apps/course/api"; name="course-api"},
    @{path="apps/learning/api"; name="learning-api"},
    @{path="apps/pay/api"; name="pay-api"},
    @{path="apps/trade/api"; name="trade-api"},
    @{path="apps/course/rpc"; name="course-rpc"},
    @{path="apps/learning/rpc"; name="learning-rpc"},
    @{path="apps/media/rpc"; name="media-rpc"},
    @{path="apps/pay/rpc"; name="pay-rpc"},
    @{path="apps/trade/rpc"; name="trade-rpc"}
)

if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

foreach ($svc in $services) {
    $fullPath = Join-Path $ProjectRoot $svc.path
    $outPath = Join-Path $ProjectRoot "bin" $svc.name
    Write-Host "  构建 $($svc.path) -> bin/$($svc.name)..."
    Set-Location $fullPath
    go build -o $outPath .
    Set-Location $ProjectRoot
}

Write-Host "构建完成！二进制文件在 bin/ 目录"