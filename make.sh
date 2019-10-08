protoc --proto_path=proto --proto_path=third_party --go_out=plugins=grpc:proto service.proto

export OLDGOOS=$GOOS
export OLDGOARCH=$GOARCH

export GOARCH=amd64
export GOOS=linux
go build -o distr/WncPlugin-linux-x64

export GOARCH=386
export GOOS=linux
go build -o distr/WncPlugin-linux-x86
+
export GOARCH=amd64
export GOOS=windows
go build -o distr/WncPlugin-windows-x64.exe

export GOARCH=amd64
export GOOS=darwin
go build -o distr/WncPlugin-macos-x64

export GOARCH=$OLDGOARCH
export GOOS=$OLDGOOS