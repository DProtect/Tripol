$ErrorActionPreference = 'Stop'
if ($env:PROCESSOR_ARCHITECTURE -ne 'AMD64') { throw 'TRIPOL requires Windows x64 PowerShell.' }
$name = 'TRIPOL-v0.28-Installer.zip'
$expected = 'd23b890e9adb8809c1a26aa7da15e79d5524ae7983dcb9dda535f8e669c288c3'
$temp = Join-Path ([IO.Path]::GetTempPath()) ('tripol-' + [Guid]::NewGuid())
New-Item -ItemType Directory -Path $temp | Out-Null
try {
    $zip = Join-Path $temp $name
    Invoke-WebRequest -UseBasicParsing -Uri ('https://github.com/Tripollas/Tripol/releases/download/v0.28.0/' + $name) -OutFile $zip
    if ((Get-FileHash -Algorithm SHA256 -LiteralPath $zip).Hash.ToLowerInvariant() -ne $expected) { throw 'SHA-256 mismatch; no installation performed.' }
    Expand-Archive -LiteralPath $zip -DestinationPath $temp
    $exe = Join-Path $temp 'TRIPOL-v0.28-Setup-x64.exe'
    $process = Start-Process -FilePath $exe -ArgumentList '--quiet' -Wait -PassThru
    if ($process.ExitCode -ne 0) { throw "Installer failed: $($process.ExitCode)" }
    Write-Host 'Open a new terminal and run tripol --version.'
} finally { Remove-Item -LiteralPath $temp -Recurse -Force }
