set GOARCH=amd64
set GOOS=linux
go build -o WncPlugin-linux-x64

set GOARCH=amd64
set GOOS=windows
go build -o WncPlugin-windows-x64.exe