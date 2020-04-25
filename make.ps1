protoc --proto_path=proto --go_out=plugins=grpc:proto commands/commands.proto
protoc --proto_path=proto --go_out=plugins=grpc:proto files/files.proto

Write-Output("Current GOOS $env:GOOS : Current GOARCH :$env:GOARCH.`n")
$env:OLDGOOS=$env:GOOS
$env:OLDGOARCH=$env:GOARCH

$env:GOARCH="amd64"
$env:GOOS="linux"
Write-Output("Linux x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -ldflags="-w -s" -o distr/WncPlugin-linux-x64
Write-Output("Linux build done.`n")

$env:GOARCH="386"
$env:GOOS="linux"
Write-Output("Linux x86 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -o distr/WncPlugin-linux-x86
Write-Output("Linux x86 build done.`n")

$env:GOARCH="amd64"
$env:GOOS="windows"
Write-Output("Windows x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -o distr/WncPlugin-windows-x64.exe
Write-Output("Windows x64 build done.`n")

$env:GOARCH="amd64"
$env:GOOS="darwin"
Write-Output("MacOS x64 build| GOOS $env:GOOS : GOARCH :$env:GOARCH.")
go build -o distr/WncPlugin-macos-x64
Write-Output("MacOS x64 build done.")

$env:GOOS=$env:OLDGOOS
$env:GOARCH=$env:OLDGOARCH