.PHONY: build build-all install test clean

VERSION := 0.1.0
BINARY := issue-flow

build:
	go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY) main.go

build-all:
	GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-darwin-amd64 main.go
	GOOS=darwin GOARCH=arm64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-darwin-arm64 main.go
	GOOS=linux GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-linux-amd64 main.go
	GOOS=windows GOARCH=amd64 go build -ldflags="-X main.version=$(VERSION)" -o bin/$(BINARY)-windows-amd64.exe main.go

install:
	go install -ldflags="-X main.version=$(VERSION)"

test:
	go test -v ./...

test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

clean:
	rm -rf bin/
	rm -f coverage.out
