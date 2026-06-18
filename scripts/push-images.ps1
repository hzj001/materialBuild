# 本地构建并推送镜像到 GitHub Container Registry
# 用法:
#   .\scripts\push-images.ps1
#   .\scripts\push-images.ps1 -Tag v1.0.0
# 需先登录: gh auth login  或  echo $env:GITHUB_TOKEN | docker login ghcr.io -u hzj001 --password-stdin

param(
    [string]$Tag = "latest",
    [string]$Registry = "ghcr.io/hzj001/materialbuild"
)

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot

function Build-Push {
    param(
        [string]$Name,
        [string]$Context,
        [string]$Dockerfile,
        [hashtable]$BuildArgs = @{}
    )
    $image = "${Registry}-${Name}:${Tag}"
    Write-Host ">>> Building $image" -ForegroundColor Cyan
    $args = @("build", "-t", $image, "-f", $Dockerfile)
    foreach ($k in $BuildArgs.Keys) {
        $args += @("--build-arg", "${k}=$($BuildArgs[$k])")
    }
    $args += $Context
    docker @args
    if ($LASTEXITCODE -ne 0) { throw "Build failed: $Name" }
    Write-Host ">>> Pushing $image" -ForegroundColor Cyan
    docker push $image
    if ($LASTEXITCODE -ne 0) { throw "Push failed: $Name" }
}

Write-Host "Registry: $Registry, Tag: $Tag" -ForegroundColor Green

# 登录 GHCR（若尚未登录）
$null = docker push "${Registry}-backend:${Tag}" 2>&1
if ($LASTEXITCODE -ne 0) {
    Write-Host "正在登录 GHCR..." -ForegroundColor Yellow
    & "C:\Program Files\GitHub CLI\gh.exe" auth token | docker login ghcr.io -u hzj001 --password-stdin
    if ($LASTEXITCODE -ne 0) {
        Write-Host "请先执行: gh auth login" -ForegroundColor Red
        exit 1
    }
}

Build-Push -Name "backend" -Context "$Root\backend" -Dockerfile "$Root\backend\Dockerfile"
Build-Push -Name "client" -Context "$Root\frontend" -Dockerfile "$Root\frontend\docker\Dockerfile" -BuildArgs @{ APP_DIR = "client" }
Build-Push -Name "merchant" -Context "$Root\frontend" -Dockerfile "$Root\frontend\docker\Dockerfile" -BuildArgs @{ APP_DIR = "merchant" }
Build-Push -Name "admin" -Context "$Root\frontend\admin" -Dockerfile "$Root\frontend\admin\Dockerfile"

Write-Host "`n全部镜像已推送:" -ForegroundColor Green
@("backend", "client", "merchant", "admin") | ForEach-Object {
    Write-Host "  ${Registry}-${_}:${Tag}"
}
