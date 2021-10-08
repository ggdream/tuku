.PHONY: all linux macos win test clean

all: win

linux:
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./cmd/tuku/tuku -ldflags '-s -w' ./cmd/tuku/main.go
macos:
	@GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -o ./cmd/tuku/tuku -ldflags '-s -w' ./cmd/tuku/main.go

win:
	@set GOOS=windows
	@set GOARCH=amd64
	@set CGO_ENABLED=0
	@go build -o ./cmd/tuku/tuku.exe -ldflags '-s -w' ./cmd/tuku/main.go


test:
	@go test -v ./...

clean:
	@del .\cmd\tuku\tuku.exe
