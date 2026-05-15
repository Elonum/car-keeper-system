# CarKeeper backend test runner
Set-Location $PSScriptRoot\..
$env:ENV = "test"
go test ./... -count=1
exit $LASTEXITCODE
