.PHONY: all linux macos win test clean

all: linux


linux:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/tuku -v -ldflags '-s -w' ./cmd/main.go

macos:
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./cmd/tuku -v -ldflags '-s -w' ./cmd/main.go

win:
	@set GOOS=windows
	@set GOARCH=amd64
	@set CGO_ENABLED=0
	@go build -o ./cmd/tuku.exe -v -ldflags '-s -w' ./cmd/main.go


test:
	@go test -v ./...

clean:
	@del .\cmd\tuku.exe
