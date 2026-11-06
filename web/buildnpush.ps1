param([Parameter(Mandatory = $true)][string]$ImageRef)
$ErrorActionPreference = 'Stop'
Push-Location $PSScriptRoot
try {
    npm.cmd run build
    if ($LASTEXITCODE -ne 0) { throw 'Frontend build failed' }
    docker build -t $ImageRef .
    if ($LASTEXITCODE -ne 0) { throw 'Docker build failed' }
    docker push $ImageRef
    if ($LASTEXITCODE -ne 0) { throw 'Docker push failed' }
} finally { Pop-Location }
