export GOARCH=amd64
export GOOS=linux
go build -o WncPlugin-linux-x64

export GOARCH=amd64
export GOOS=windows
go build -o WncPlugin-windows-x64.exe