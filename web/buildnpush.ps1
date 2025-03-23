$ErrorActionPreference = "Stop"

$IMAGE_NAME = "maasst"
$REGISTRY = "registry.cn-chengdu.aliyuncs.com/example/maa-web"
$TAG = "latest"

try {
    npm run build
    docker build -t ${IMAGE_NAME} .
    docker tag ${IMAGE_NAME} "${REGISTRY}:${TAG}"
    docker push "${REGISTRY}:${TAG}"
    Write-Host "`n[SUCCESS] All operations completed!" -ForegroundColor Green
}
catch {
    Write-Host "`n[ERROR] Failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}