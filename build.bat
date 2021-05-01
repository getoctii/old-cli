mkdir builds

set GOOS=windows
set GOARCH=amd64
go build -o builds/win64.exe main.go

set GOOS=windows
set GOARCH=386
go build -o builds/win32.exe main.go

set GOOS=darwin
set GOARCH=amd64
go build -o builds/intel-macos main.go

set GOOS=darwin
set GOARCH=arm64
go build -o builds/arm-macos main.go

set GOOS=linux
set GOARCH=amd64
go build -o builds/amd64-linux main.go

set GOOS=linux
set GOARCH=arm64
go build -o builds/arm64-linux main.go