# 建材通 Docker 开发环境启动脚本 (Windows PowerShell)
param(
    [switch]$Detached
)

$ErrorActionPreference = "Stop"
Set-Location $PSScriptRoot\..

if (-not (Test-Path ".env")) {
    Copy-Item ".env.example" ".env"
    Write-Host "已创建 .env 文件，可按需修改配置" -ForegroundColor Yellow
}

$args = @("-f", "docker-compose.yml", "-f", "docker-compose.dev.yml", "up", "--build")
if ($Detached) { $args += "-d" }

Write-Host "启动 materialBuild 开发环境..." -ForegroundColor Green
docker compose @args

if ($Detached) {
    Write-Host ""
    Write-Host "服务地址:" -ForegroundColor Cyan
    Write-Host "  客户端:   http://localhost:4000"
    Write-Host "  商家端:   http://localhost:4001"
    Write-Host "  管理后台: http://localhost:4002  (vue-pure-admin)"
    Write-Host "  后端 API: http://localhost:4008/health"
    Write-Host "  MySQL:    localhost:4306"
    Write-Host "  Redis:    localhost:4379"
}
