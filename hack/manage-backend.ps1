param(
    [ValidateSet('Status', 'Restart', 'Update')][string]$Action = 'Status',
    [string]$Candidate,
    [switch]$Build,
    [string]$BaseUrl = 'http://127.0.0.1:5244',
    [string]$ServiceName = 'maa-backend'
)
$ErrorActionPreference = 'Stop'
$backendRoot = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..'))
$installed = Join-Path $backendRoot 'main.exe'
$maintenanceDir = Join-Path $backendRoot ('.runtime\maintenance\' + (Get-Date -Format 'yyyyMMdd-HHmmss-fff'))
$node = (Get-Command node -ErrorAction Stop).Source
$healthScript = Join-Path $PSScriptRoot 'backend-health.cjs'
Set-Location -LiteralPath $backendRoot
if ($Build) {
    if ($Action -ne 'Update') { throw '-Build requires -Action Update' }
    [IO.Directory]::CreateDirectory($maintenanceDir) | Out-Null
    $Candidate = Join-Path $maintenanceDir 'candidate.exe'
    & go build -o $Candidate .
    if ($LASTEXITCODE -ne 0) { throw 'Build failed; backend was not changed.' }
}
if ($Candidate) { $Candidate = (Resolve-Path -LiteralPath $Candidate).Path }
if ($Action -eq 'Update' -and -not $Candidate) { throw 'Specify -Candidate <exe> or -Build.' }
if ($Action -eq 'Update' -and $Candidate -eq $installed) { throw 'Build a separate candidate instead of overwriting running main.exe.' }
if ($Action -eq 'Status') {
    [IO.Directory]::CreateDirectory($maintenanceDir) | Out-Null
    & $node $healthScript snapshot --snapshot (Join-Path $maintenanceDir 'status.json') --base-url $BaseUrl
    exit $LASTEXITCODE
}
$identity = [Security.Principal.WindowsIdentity]::GetCurrent()
$principal = New-Object Security.Principal.WindowsPrincipal($identity)
if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
    $arguments = '-NoProfile -ExecutionPolicy Bypass -File "' + $PSCommandPath + '" -Action ' + $Action + ' -BaseUrl "' + $BaseUrl + '" -ServiceName "' + $ServiceName + '"'
    if ($Candidate) { $arguments += ' -Candidate "' + $Candidate + '"' }
    $elevated = Start-Process -FilePath 'powershell.exe' -Verb RunAs -WindowStyle Hidden -ArgumentList $arguments -Wait -PassThru
    exit $elevated.ExitCode
}
$mutex = New-Object System.Threading.Mutex($false, 'Global\CodexMaaMuMuDeployment')
if (-not $mutex.WaitOne(0)) { throw 'Another backend deployment is already running.' }
$stopped = $false
$replaced = $false
$backup = Join-Path $maintenanceDir 'main-before.exe'
$resultFile = Join-Path $maintenanceDir 'result.json'
$stage = 'preflight'
function Stop-BackendService {
    Stop-Service -Name $ServiceName
    (Get-Service -Name $ServiceName).WaitForStatus('Stopped', [TimeSpan]::FromSeconds(45))
    $remaining = Get-Process -Id $oldProcess.Id -ErrorAction SilentlyContinue
    if ($remaining) {
        if ($remaining.Path -ne $installed -or $remaining.SessionId -ne 0 -or $remaining.StartTime -ne $oldStart) { throw 'Unexpected surviving process' }
        Stop-Process -Id $remaining.Id -Force
        if (-not $remaining.WaitForExit(10000)) { throw 'Backend did not exit' }
    }
}
try {
    [IO.Directory]::CreateDirectory($maintenanceDir) | Out-Null
    $service = Get-Service -Name $ServiceName
    if ($service.Status -ne 'Running') { throw 'Service is not running; inspect it before updating.' }
    $oldListener = Get-NetTCPConnection -LocalPort ([uri]$BaseUrl).Port -State Listen | Select-Object -First 1
    $oldProcess = Get-Process -Id $oldListener.OwningProcess
    if ($oldProcess.Path -ne $installed -or $oldProcess.SessionId -ne 0) { throw 'Unexpected backend executable or session.' }
    $oldStart = $oldProcess.StartTime
    $snapshot = Join-Path $maintenanceDir 'before.json'
    & $node $healthScript snapshot --snapshot $snapshot --base-url $BaseUrl
    if ($LASTEXITCODE -ne 0) { throw 'Accounts are busy or health checks failed; backend was not stopped.' }
    Copy-Item -LiteralPath $installed -Destination $backup
    if ($Action -eq 'Update') { $candidateHash = (Get-FileHash -LiteralPath $Candidate).Hash }
    $stage = 'stop-service'
    $stopped = $true
    Stop-BackendService
    if ($Action -eq 'Update') {
        $stage = 'replace-executable'
        $replaced = $true
        Copy-Item -LiteralPath $Candidate -Destination $installed -Force
        if ((Get-FileHash -LiteralPath $installed).Hash -ne $candidateHash) { throw 'Installed hash mismatch' }
    }
    $stage = 'restart-and-verify'
    Start-Service -Name $ServiceName
    (Get-Service -Name $ServiceName).WaitForStatus('Running', [TimeSpan]::FromSeconds(30))
    & $node $healthScript verify --snapshot $snapshot --base-url $BaseUrl
    if ($LASTEXITCODE -ne 0) { throw 'New backend health check failed' }
    $listener = Get-NetTCPConnection -LocalPort ([uri]$BaseUrl).Port -State Listen | Select-Object -First 1
    $process = Get-Process -Id $listener.OwningProcess
    if ($process.Path -ne $installed -or $process.SessionId -ne 0 -or $process.Id -eq $oldProcess.Id) { throw 'Unexpected restarted process' }
    @{success=$true; action=$Action; service=$ServiceName; processId=$process.Id; sessionId=$process.SessionId; backup=$backup; sha256=(Get-FileHash -LiteralPath $installed).Hash; time=(Get-Date -Format o)} | ConvertTo-Json | Set-Content -LiteralPath $resultFile -Encoding UTF8
    Write-Output ('Backend ' + $Action + ' completed; sessions and schedules verified. Result: ' + $resultFile)
} catch {
    $rollbackSucceeded = $null
    if ($stopped) {
        try {
            Stop-BackendService
            if ($replaced) { Copy-Item -LiteralPath $backup -Destination $installed -Force }
            Start-Service -Name $ServiceName
            (Get-Service -Name $ServiceName).WaitForStatus('Running', [TimeSpan]::FromSeconds(30))
            & $node $healthScript verify --snapshot $snapshot --base-url $BaseUrl
            if ($LASTEXITCODE -ne 0) { throw 'Rollback health check failed' }
            $rollbackSucceeded = $true
        } catch { $rollbackSucceeded = $false }
    }
    @{success=$false; action=$Action; stage=$stage; rollbackSucceeded=$rollbackSucceeded; time=(Get-Date -Format o)} | ConvertTo-Json | Set-Content -LiteralPath $resultFile -Encoding UTF8
    Write-Error ('Backend maintenance failed during ' + $stage + '; see ' + $resultFile + '. Rollback: ' + $rollbackSucceeded) -ErrorAction Continue
    exit 1
} finally {
    $mutex.ReleaseMutex()
    $mutex.Dispose()
}
