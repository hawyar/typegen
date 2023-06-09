.PHONY: test build

name = typegen
os = $(shell go env GOOS)
arch = $(shell go env GOARCH)

test:
	go test -v ./...

build:
	GOOS=$(os)  GOARCH=$(arch) go build -o $(name) .
