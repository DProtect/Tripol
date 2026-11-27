$ErrorActionPreference = 'Stop'
if ($env:PROCESSOR_ARCHITECTURE -ne 'AMD64') { throw 'TRIPOL requires Windows x64 PowerShell.' }
$name = 'TRIPOL-v0.27-Installer.zip'
$expected = 'e5ae51376b4ed68da05dcc88c890fb86b9b1e1161ec285a2a2cab732a471ce1c'
$temp = Join-Path ([IO.Path]::GetTempPath()) ('tripol-' + [Guid]::NewGuid())
New-Item -ItemType Directory -Path $temp | Out-Null
try {
    $zip = Join-Path $temp $name
    $local = Join-Path $PSScriptRoot ('downloads/' + $name)
    if (Test-Path -LiteralPath $local) { Copy-Item -LiteralPath $local -Destination $zip }
    else { Invoke-WebRequest -UseBasicParsing -Uri ('https://github.com/Tripollas/Tripol/releases/download/v0.27.0/' + $name) -OutFile $zip }
    if ((Get-FileHash -Algorithm SHA256 -LiteralPath $zip).Hash.ToLowerInvariant() -ne $expected) { throw 'SHA-256 mismatch; no installation performed.' }
    Expand-Archive -LiteralPath $zip -DestinationPath $temp
    $exe = Join-Path $temp 'TRIPOL-v0.27-Setup-x64.exe'
    $process = Start-Process -FilePath $exe -ArgumentList '--quiet' -Wait -PassThru
    if ($process.ExitCode -ne 0) { throw "Installer failed: $($process.ExitCode)" }
    Write-Host 'Open a new terminal and run tripol --version.'
} finally { Remove-Item -LiteralPath $temp -Recurse -Force }
