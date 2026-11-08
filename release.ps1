param(
  [ValidateSet("windows","linux","darwin")][string]$OS,
  [string]$Version,
  [switch]$arm #arm64
)

if (-not $OS -or -not $Version) {
  Write-Host "Usage: .\release.ps1 <windows|linux|darwin> <version> [ -arm ]"
  exit 1
}

$GOARCH = "amd64"
if ($arm) { $GOARCH = "arm64" }

$EXT = ""
if ($OS -ieq "windows") { $EXT = ".exe" }

# set env for this process only
$env:GOOS = $OS
$env:GOARCH = $GOARCH
$env:CGO_ENABLED = "0"
$OutName = "bench-cmp-$($env:GOOS)-$($env:GOARCH)-$Version$EXT"

Write-Host "Building $OutName (GOOS=$env:GOOS GOARCH=$env:GOARCH CGO_ENABLED=$env:CGO_ENABLED)"

$build = Start-Process go -ArgumentList "build","-o",$OutName,"./cmd" -NoNewWindow -PassThru -Wait
if ($build.ExitCode -ne 0) {
  Write-Error "Build failed (exit $($build.ExitCode))"
  exit $build.ExitCode
}

Write-Host "Created $OutName"
