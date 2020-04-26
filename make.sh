protoc --proto_path=proto --go_out=plugins=grpc:proto commands/commands.proto
protoc --proto_path=proto --go_out=plugins=grpc:proto files/files.proto

export PRODUCT_NAME="PLMCompanionAddon"
echo "Building $PRODUCT_NAME"

echo "Current GOOS $GOOS : Current GOARCH $GOARCH"

export GOARCH=amd64
export GOOS=linux
echo "Linux x64 build| GOOS $GOOS : GOARCH :$GOARCH."
go build -o distr/$PRODUCT_NAME-linux-x64
echo "Linux build done."

export GOARCH=386
export GOOS=linux
echo "Linux x86 build| GOOS $GOOS : GOARCH :$GOARCH."
go build -o distr/$PRODUCT_NAME-linux-x86
echo "Linux x86 build done."

export GOARCH=amd64
export GOOS=windows
echo "Windows x64 build| GOOS $GOOS : GOARCH :$GOARCH."
go build -o distr/$PRODUCT_NAME-windows-x64.exe
echo "Windows x64 build done."

export GOARCH=amd64
export GOOS=darwin
echo "MacOS x64 build| GOOS $GOOS : GOARCH :$GOARCH."
go build -o distr/$PRODUCT_NAME-macos-x64
echo "MacOS x64 build done."

export GOARCH=$OLDGOARCH
export GOOS=$OLDGOOS