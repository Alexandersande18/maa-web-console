param([string]$cert = 'false', [string]$ImageRef = '')
$ErrorActionPreference = 'Stop'
Push-Location $PSScriptRoot
try {
    $useCert = ($cert -replace '^cert=', '').ToLowerInvariant() -eq 'true'
    $localConfigPath = Join-Path $PSScriptRoot 'deploy.local.json'
    if (-not $ImageRef -and (Test-Path -LiteralPath $localConfigPath)) {
        $localConfig = Get-Content -LiteralPath $localConfigPath -Raw | ConvertFrom-Json
        if ($useCert) { $ImageRef = $localConfig.certImageRef } else { $ImageRef = $localConfig.imageRef }
    }
    if (-not $ImageRef) { throw 'Pass -ImageRef or configure the ignored deploy.local.json file.' }
    $dockerfile = if ($useCert) { 'withcert.dockerfile' } else { 'Dockerfile' }
    npm.cmd run build
    if ($LASTEXITCODE -ne 0) { throw 'Frontend build failed' }
    docker build -f $dockerfile -t $ImageRef .
    if ($LASTEXITCODE -ne 0) { throw 'Docker build failed' }
    docker push $ImageRef
    if ($LASTEXITCODE -ne 0) { throw 'Docker push failed' }
} finally { Pop-Location }
