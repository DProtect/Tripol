param(
    [string]$Compiler = 'tripol',
    [switch]$Benchmark,
    [switch]$CheckUpdates
)
$ErrorActionPreference = 'Stop'
$destination = Join-Path (Get-Location) 'tripol-validation-v027.json'
$benchmarkPath = Join-Path (Get-Location) 'arena-v027.json'
$root = Join-Path ([System.IO.Path]::GetTempPath()) ('tripol-v027-' + [Guid]::NewGuid().ToString('N'))
$results = [System.Collections.Generic.List[object]]::new()
$oldOpt = $env:TRIPOL_NO_UPDATE
$version = ''
$resolved = ''
$hardware = @{}
function Test-Case([string]$Name, [scriptblock]$Body) {
    $watch = [Diagnostics.Stopwatch]::StartNew()
    try {
        $detail = & $Body
        $status = 'PASS'
    } catch {
        $detail = $_.Exception.Message
        $status = 'FAIL'
    }
    $watch.Stop()
    $results.Add([pscustomobject]@{name=$Name;status=$status;duration_ms=$watch.ElapsedMilliseconds;detail=($detail -join "`n")})
    Write-Host "[$status] $Name"
}
function Invoke-Tripol([string[]]$Arguments) {
    $output = & $resolved @Arguments 2>&1
    if ($LASTEXITCODE -ne 0) { throw ($output -join "`n") }
    return ($output -join "`n")
}
function Assert-Output([string]$Actual,[string]$Expected) {
    if (($Actual -replace "`r","").Trim() -ne $Expected.Trim()) { throw "Expected '$Expected'; got '$Actual'" }
}
Write-Host 'TRIPOL QUANTUM X - VALIDATION CENTER'
try {
    $env:TRIPOL_NO_UPDATE = '1'
    New-Item -ItemType Directory -Path $root | Out-Null
    $resolved = (Get-Command $Compiler -CommandType Application).Source
    $version = Invoke-Tripol @('--version')
    try {
        $cpu = Get-CimInstance Win32_Processor | Select-Object -First 1
        $os = Get-CimInstance Win32_OperatingSystem
        $hardware = @{cpu=$cpu.Name;ram_bytes=($os.TotalVisibleMemorySize * 1024);os=$os.Caption;arch=$env:PROCESSOR_ARCHITECTURE}
    } catch { $hardware = @{arch=$env:PROCESSOR_ARCHITECTURE} }
    Test-Case 'Installation' { Assert-Output $version 'TRIPOL 0.27.1'; $version }
    $source = Join-Path $root 'main.tri'
    @'
import tripol.math
a=ints(2)
a[0]=7
a[0]+=1
show math.square(a[0])
try {a[0]/=0} catch e {show e}
try {a[2]+=1} catch e {show e}
'@ | Set-Content -LiteralPath $source -Encoding ASCII
    Test-Case 'Compiler check' { Invoke-Tripol @('check',$source) }
    foreach ($profile in @('nano','balanced','turbo','debug')) {
        Test-Case "Native $profile" {
            $exe = Join-Path $root "$profile.exe"
            Invoke-Tripol @('build',$source,'--target','windows','--profile',$profile,'-o',$exe) | Out-Null
            $actual = & $exe
            if ($LASTEXITCODE -ne 0) { throw "Native program exited $LASTEXITCODE" }
            Assert-Output ($actual -join "`n") "64`n10`n11"
            "Bytes: $((Get-Item -LiteralPath $exe).Length)"
        }
    }
    Test-Case 'Invalid source rejection' {
        $bad = Join-Path $root 'invalid.tri'
        'show unknown_symbol' | Set-Content -LiteralPath $bad -Encoding ASCII
        $p = Start-Process -FilePath $resolved -ArgumentList @('check',('"' + $bad + '"')) -Wait -PassThru -NoNewWindow -RedirectStandardError (Join-Path $root 'error.txt') -RedirectStandardOutput (Join-Path $root 'output.txt')
        if ($p.ExitCode -eq 0) { throw 'Invalid source was accepted' }
        Get-Content -LiteralPath (Join-Path $root 'error.txt')
    }
    Test-Case 'New libraries' {
        $lib = Join-Path $root 'libraries.tri'
        @'
import tripol.heap
import tripol.utf8
import tripol.csv
h=heap.create(4)
assert heap.push(h,9)
assert heap.push(h,2)
show heap.pop(h)
show utf8.count("abc")
show csv.fields("a,b,c")
'@ | Set-Content -LiteralPath $lib -Encoding ASCII
        $exe = Join-Path $root 'libraries.exe'
        Invoke-Tripol @('build',$lib,'--target','windows','-o',$exe) | Out-Null
        $actual = & $exe
        if ($LASTEXITCODE -ne 0) { throw "Native program exited $LASTEXITCODE" }
        Assert-Output ($actual -join "`n") "2`n3`n3"
    }
    Test-Case 'Automatic loop memory reclamation' {
        $memory = Join-Path $root 'memory.tri'
        @'
import tripol.memory
before=memory.used()
repeat 1000 {a=ints(65535);a[65534]=9}
show memory.used()==before
'@ | Set-Content -LiteralPath $memory -Encoding ASCII
        $exe = Join-Path $root 'memory.exe'
        Invoke-Tripol @('build',$memory,'--target','windows','-o',$exe) | Out-Null
        $actual = & $exe
        if ($LASTEXITCODE -ne 0) { throw "Native program exited $LASTEXITCODE" }
        Assert-Output ($actual -join "`n") 'true'
    }
    Test-Case 'Dynamic lists and aliases' {
        $list = Join-Path $root 'list.tri'
        @'
a=list()
b=a
repeat 100 {a.add(7)}
show b.size(),b.get(99)
a.remove_at(0)
show b.size()
'@ | Set-Content -LiteralPath $list -Encoding ASCII
        $exe = Join-Path $root 'list.exe'
        Invoke-Tripol @('build',$list,'--target','windows','-o',$exe) | Out-Null
        $actual = & $exe
        if ($LASTEXITCODE -ne 0) { throw "Native program exited $LASTEXITCODE" }
        Assert-Output ($actual -join "`n") "100 7`n99"
    }
    if ($CheckUpdates) { Test-Case 'Live update detection' { Invoke-Tripol @('update','--check') } }
    else { $results.Add([pscustomobject]@{name='Live update detection';status='SKIP';detail='Use -CheckUpdates'}) }
    if ($Benchmark) { Test-Case 'Arena' { Invoke-Tripol @('arena','--runs','20','--profile','turbo','--tripol-only','--json',$benchmarkPath) } }
    else { $results.Add([pscustomobject]@{name='Arena';status='SKIP';detail='Use -Benchmark'}) }
} catch {
    $results.Add([pscustomobject]@{name='Setup';status='FAIL';detail=$_.Exception.Message})
    Write-Host "[FAIL] Setup: $($_.Exception.Message)"
} finally {
    $env:TRIPOL_NO_UPDATE = $oldOpt
    if (Test-Path -LiteralPath $root) { Remove-Item -LiteralPath $root -Recurse -Force }
}
$passed = @($results | Where-Object status -eq 'PASS').Count
$failed = @($results | Where-Object status -eq 'FAIL').Count
$skipped = @($results | Where-Object status -eq 'SKIP').Count
[pscustomobject]@{version=$version;hardware=$hardware;passed=$passed;failed=$failed;skipped=$skipped;tests=$results} | ConvertTo-Json -Depth 7 | Set-Content -LiteralPath $destination -Encoding UTF8
Write-Host "Passed: $passed; Failed: $failed; Skipped: $skipped"
Write-Host "Report: $destination"
if ($failed -gt 0) { exit 1 }
