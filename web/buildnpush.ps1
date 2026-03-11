<#  buildnpush.ps1
    调用示例：
        .\buildnpush.ps1             # 用默认 Dockerfile 和 latest tag
        .\buildnpush.ps1 cert=true   # 用 withcert.dockerfile 和 cert tag
#>

param(
    [string]$cert = "false"          # 允许 PowerShell 7 里的 --cert=true 形式
)

$ErrorActionPreference = "Stop"

# ─── 基本变量 ─────────────────────────────────────────────
$IMAGE_NAME = "maasst"
$REGISTRY   = "registry.cn-chengdu.aliyuncs.com/example/maa-web"

# ─── 选择是否使用证书 ───────────────────────────────────────
$useCert = $false
if ($cert.ToLower() -eq "true") {
    $useCert = $true
}

# ─── 根据是否使用证书设置 Dockerfile 和 Tag ──────────────
if ($useCert) {
    $dockerfile = "withcert.dockerfile"
    $TAG        = "cert"
} else {
    $dockerfile = "Dockerfile"
    $TAG        = "latest"
}

Write-Host "`n[INFO] Using Dockerfile: $dockerfile" -ForegroundColor Cyan
Write-Host "[INFO] Using Tag: $TAG" -ForegroundColor Cyan

# 将完整镜像名单独存到变量，避免冒号解析问题
$FULL_TAG = "${REGISTRY}:${TAG}"

# ─── 主流程 ───────────────────────────────────────────────
try {
    Write-Host "[STEP] Front-end build (npm run build) ..." -ForegroundColor Yellow
    npm run build

    Write-Host "`n[STEP] Docker build ..." -ForegroundColor Yellow
    docker build -f $dockerfile -t $IMAGE_NAME .

    Write-Host "`n[STEP] Docker tag ..." -ForegroundColor Yellow
    docker tag $IMAGE_NAME $FULL_TAG

    Write-Host "`n[STEP] Docker push ..." -ForegroundColor Yellow
    docker push $FULL_TAG

    Write-Host "`n[SUCCESS] All operations completed!" -ForegroundColor Green
}
catch {
    Write-Host "`n[ERROR] Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}