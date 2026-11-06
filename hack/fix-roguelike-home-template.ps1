param(
    [Parameter(Mandatory = $true)]
    [string]$ResourcePath
)

$ErrorActionPreference = 'Stop'
$root = [IO.Path]::GetFullPath($ResourcePath)
$taskFile = Join-Path $root 'resource\tasks\Roguelike\base.json'
$raw = [IO.File]::ReadAllText($taskFile)
$task = ($raw | ConvertFrom-Json).Roguelike
$templates = @($task.template)
if ($templates -contains 'SwitchTheme@ToggleSettingsMenu-v2.png') {
    Write-Output "Already fixed: $taskFile"
    return
}
if ($templates.Count -ne 1 -or $templates[0] -ne 'SwitchTheme@ToggleSettingsMenu.png') {
    throw 'Unexpected Roguelike template; refusing to overwrite this resource version'
}
$templateFile = Join-Path $root 'resource\template\UiTheme\SwitchTheme\SwitchTheme@ToggleSettingsMenu-v2.png'
if (-not (Test-Path -LiteralPath $templateFile -PathType Leaf)) {
    throw "Missing v2 home template: $templateFile"
}

# Same one-field change as upstream PR #19010. Preserve all other task data.
$pattern = '(?s)("Roguelike"\s*:\s*\{[^{}]*?"template"\s*:\s*)"SwitchTheme@ToggleSettingsMenu\.png"'
$matcher = [regex]::new($pattern)
if ($matcher.Matches($raw).Count -ne 1) {
    throw 'Could not uniquely locate the Roguelike template field'
}
$patched = $matcher.Replace($raw, '${1}["SwitchTheme@ToggleSettingsMenu.png", "SwitchTheme@ToggleSettingsMenu-v2.png"]', 1)
$null = $patched | ConvertFrom-Json
$backupDir = Join-Path $PSScriptRoot '..\.runtime\resource-backups'
$null = New-Item -ItemType Directory -Path $backupDir -Force
$backup = Join-Path $backupDir ('roguelike-base-' + (Get-Date -Format 'yyyyMMdd-HHmmss-fff') + '.json')
$temporary = $taskFile + '.' + [Guid]::NewGuid().ToString('N') + '.tmp'
[IO.File]::WriteAllText($temporary, $patched, [Text.UTF8Encoding]::new($false))
[IO.File]::Replace($temporary, $taskFile, $backup)
Write-Output "Patched: $taskFile"
Write-Output "Backup: $([IO.Path]::GetFullPath($backup))"
Write-Output 'The backend reloads this resource before the next queued batch.'
