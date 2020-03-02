protoc --proto_path=proto --go_out=plugins=grpc:proto commands/commands.proto
protoc --proto_path=proto --go_out=plugins=grpc:proto files/files.proto

set OLDGOOS=%GOOS%
set OLDGOARCH=%GOARCH%

set GOARCH=amd64
set GOOS=linux
go build -o distr/WncPlugin-linux-x64

set GOARCH=386
set GOOS=linux
go build -o distr/WncPlugin-linux-x86

set GOARCH=amd64
set GOOS=windows
go build -o distr/WncPlugin-windows-x64.exe

set GOARCH=amd64
set GOOS=darwin
go build -o distr/WncPlugin-macos-x64

set GOOS=%OLDGOOS%
set GOARCH=%OLDGOARCH%