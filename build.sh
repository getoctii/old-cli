mkdir builds
GOOS="windows" GOARCH="amd64" go build -o builds/win64.exe main.go
GOOS="windows" GOARCH="386" go build -o builds/win32.exe main.go
GOOS="darwin" GOARCH="amd64" go build -o builds/intel-macos main.go
GOOS="darwin" GOARCH="arm64" go build -o builds/arm-macos main.go
GOOS="linux" GOARCH="amd64" go build -o builds/amd64-linux main.go
GOOS="linux" GOARCH="arm64" go build -o builds/arm64-linux main.go