# materialBuild full integration test
param(
    [string]$BaseUrl = "http://localhost:4008",
    [string]$ClientUrl = "http://localhost:4000",
    [string]$MerchantUrl = "http://localhost:4001",
    [string]$AdminUrl = "http://localhost:4002"
)

$ErrorActionPreference = "Continue"
$passed = 0
$failed = 0
$results = @()

function Test-Case {
    param([string]$Name, [scriptblock]$Block)
    try {
        & $Block
        $script:passed++
        $script:results += [PSCustomObject]@{ Name = $Name; Status = "PASS"; Detail = "" }
        Write-Host "[PASS] $Name" -ForegroundColor Green
    } catch {
        $script:failed++
        $msg = $_.Exception.Message
        $script:results += [PSCustomObject]@{ Name = $Name; Status = "FAIL"; Detail = $msg }
        Write-Host "[FAIL] $Name - $msg" -ForegroundColor Red
    }
}

Write-Host ""
Write-Host "========== materialBuild Test Report ==========" -ForegroundColor Cyan

Test-Case "Backend /health" {
    $r = Invoke-RestMethod -Uri "$BaseUrl/health" -TimeoutSec 10
    if ($r.status -ne "ok") { throw "status not ok" }
}

Test-Case "API common/cities" {
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/common/cities" -TimeoutSec 10
    if ($r.code -ne 0 -or $r.data.Count -lt 1) { throw "empty cities" }
}

Test-Case "API common/categories" {
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/common/categories" -TimeoutSec 10
    if ($r.code -ne 0 -or $r.data.Count -lt 1) { throw "empty categories" }
}

$adminToken = $null
Test-Case "Admin login" {
    $body = '{"phone":"13800000000","password":"admin123","role":"admin"}'
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $body -ContentType "application/json; charset=utf-8" -TimeoutSec 10
    if ($r.code -ne 0 -or -not $r.data.token) { throw $r.message }
    $script:adminToken = $r.data.token
}

Test-Case "Admin stats" {
    $h = @{ Authorization = "Bearer $adminToken" }
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/stats" -Headers $h -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

Test-Case "Admin cities" {
    $h = @{ Authorization = "Bearer $adminToken" }
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/admin/cities" -Headers $h -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

Test-Case "Admin merchants" {
    $h = @{ Authorization = "Bearer $adminToken" }
    $uri = "$BaseUrl/api/v1/admin/merchants" + "?page_size=10"
    $r = Invoke-RestMethod -Uri $uri -Headers $h -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

Test-Case "Admin orders" {
    $h = @{ Authorization = "Bearer $adminToken" }
    $uri = "$BaseUrl/api/v1/admin/orders" + "?page_size=10"
    $r = Invoke-RestMethod -Uri $uri -Headers $h -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

$userToken = $null
$testPhone = "139" + (Get-Random -Minimum 10000000 -Maximum 99999999)
Test-Case "User register" {
    $body = "{`"phone`":`"$testPhone`",`"password`":`"test123456`",`"nickname`":`"testuser`",`"role`":`"user`",`"city_id`":1}"
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/login" -Method POST -Body $body -ContentType "application/json" -TimeoutSec 10 -ErrorAction SilentlyContinue
    $bodyReg = "{`"phone`":`"$testPhone`",`"password`":`"test123456`",`"nickname`":`"testuser`",`"role`":`"user`",`"city_id`":1}"
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/register" -Method POST -Body $bodyReg -ContentType "application/json" -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
    $script:userToken = $r.data.token
}

Test-Case "Client products" {
    $uri = "$BaseUrl/api/v1/client/products" + "?city_id=1" + "&page_size=5"
    $r = Invoke-RestMethod -Uri $uri -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

Test-Case "Client merchants" {
    $uri = "$BaseUrl/api/v1/client/merchants" + "?city_id=1"
    $r = Invoke-RestMethod -Uri $uri -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

$merchantPhone = "137" + (Get-Random -Minimum 10000000 -Maximum 99999999)
Test-Case "Merchant register" {
    $body = "{`"phone`":`"$merchantPhone`",`"password`":`"merchant123`",`"nickname`":`"testmerchant`",`"role`":`"merchant`"}"
    $r = Invoke-RestMethod -Uri "$BaseUrl/api/v1/auth/register" -Method POST -Body $body -ContentType "application/json" -TimeoutSec 10
    if ($r.code -ne 0) { throw $r.message }
}

Test-Case "Client H5 index" {
    $r = Invoke-WebRequest -Uri $ClientUrl -TimeoutSec 10 -UseBasicParsing
    if ($r.StatusCode -ne 200) { throw "HTTP $($r.StatusCode)" }
}

Test-Case "Client shared CSS" {
    $r = Invoke-WebRequest -Uri ($ClientUrl + "/shared/css/base.css") -TimeoutSec 10 -UseBasicParsing
    if ($r.StatusCode -ne 200 -or $r.Content -notmatch "--brand") { throw "CSS missing" }
}

Test-Case "Merchant H5 index" {
    $r = Invoke-WebRequest -Uri $MerchantUrl -TimeoutSec 10 -UseBasicParsing
    if ($r.StatusCode -ne 200) { throw "HTTP $($r.StatusCode)" }
}

Test-Case "Client nginx API proxy" {
    $uri = $ClientUrl + "/api/v1/common/cities"
    $r = Invoke-RestMethod -Uri $uri -TimeoutSec 10
    if ($r.code -ne 0) { throw "proxy failed" }
}

Test-Case "Admin panel (optional)" {
    try {
        $r = Invoke-WebRequest -Uri $AdminUrl -TimeoutSec 30 -UseBasicParsing
        if ($r.StatusCode -ne 200) { throw "HTTP $($r.StatusCode)" }
    } catch {
        throw "Admin not ready yet: $($_.Exception.Message)"
    }
}

Write-Host ""
Write-Host "========== Summary ==========" -ForegroundColor Cyan
Write-Host "Passed: $passed  Failed: $failed  Total: $($passed + $failed)"
$results | Format-Table -AutoSize
if ($failed -gt 0) { exit 1 }
