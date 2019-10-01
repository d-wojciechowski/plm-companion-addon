export GOARCH=amd64
export GOOS=linux
go build -o WncPlugin-linux-x64

export GOARCH=386
export GOOS=linux
go build -o WncPlugin-linux-x86

export GOARCH=amd64
export GOOS=windows
go build -o WncPlugin-windows-x64.exe

export GOARCH=amd64
export GOOS=darwin
go build -o WncPlugin-macos-x64.exe