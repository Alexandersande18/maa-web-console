param(
    [string]$ResourcePath = 'C:\MAA',
    [string]$AdbPath = 'C:\platform-tools\adb.exe',
    [string]$AdbAddress = '127.0.0.1:16384'
)
$ErrorActionPreference = 'Stop'
$configDir = [IO.Path]::GetFullPath((Join-Path $PSScriptRoot '..\manifest\config'))
$destination = Join-Path $configDir 'config.yaml'
if (Test-Path -LiteralPath $destination) {
    throw 'config.yaml already exists. Update it manually; existing configuration is never overwritten.'
}
function New-RandomSecret {
    $bytes = New-Object byte[] 32
    $rng = [Security.Cryptography.RandomNumberGenerator]::Create()
    try { $rng.GetBytes($bytes) } finally { $rng.Dispose() }
    return ([BitConverter]::ToString($bytes)).Replace('-', '').ToLowerInvariant()
}
$raw = [IO.File]::ReadAllText((Join-Path $configDir 'config.example.yaml'))
$raw = $raw.Replace('jwtsecret: ""', ('jwtsecret: "' + (New-RandomSecret) + '"'))
$raw = $raw.Replace('token: ""', ('token: "' + (New-RandomSecret) + '"'))
$raw = $raw.Replace('C:\MAA', $ResourcePath.Replace("'", "''"))
$raw = $raw.Replace('C:\platform-tools\adb.exe', $AdbPath.Replace("'", "''"))
if ($AdbAddress -notmatch '^[a-zA-Z0-9.\[\]:_-]+$') { throw 'Invalid ADB address' }
$raw = $raw.Replace('127.0.0.1:16384', $AdbAddress)
[IO.File]::WriteAllText($destination, $raw, [Text.UTF8Encoding]::new($false))
Write-Output 'Created manifest/config/config.yaml with independent random credentials.'
Write-Output 'Set your database connection, verify ADB settings, and read auth.users[0].token locally to log in.'
