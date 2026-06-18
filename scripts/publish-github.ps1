# 一键发布：推送代码 + 容器镜像到 GitHub
# 前置：gh auth login 已完成
# 用法: .\scripts\publish-github.ps1

$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $PSScriptRoot
Set-Location $Root

$gh = "C:\Program Files\GitHub CLI\gh.exe"
& $gh auth status 2>$null
if ($LASTEXITCODE -ne 0) {
    Write-Host "请先登录 GitHub: gh auth login" -ForegroundColor Red
    exit 1
}

Write-Host ">>> 推送代码到 GitHub..." -ForegroundColor Cyan
git push -u origin main
if ($LASTEXITCODE -ne 0) { throw "git push 失败" }

Write-Host ">>> 登录 GHCR..." -ForegroundColor Cyan
& $gh auth token | docker login ghcr.io -u hzj001 --password-stdin
if ($LASTEXITCODE -ne 0) { throw "GHCR 登录失败" }

Write-Host ">>> 推送容器镜像..." -ForegroundColor Cyan
& "$PSScriptRoot\push-images.ps1"

Write-Host "`n发布完成!" -ForegroundColor Green
Write-Host "代码: https://github.com/hzj001/materialBuild"
Write-Host "镜像: ghcr.io/hzj001/materialbuild-{backend|client|merchant|admin}:latest"
Write-Host "GitHub Actions 将在 push 后自动构建镜像（与本地推送并行可用）"
