.PHONY: all build race build-download_channels build-download build-search

all: lint test race build build-download_channels build-download build-search

init:
	go mod tidy

fmt:
	go fmt ./...

lint:
	go vet ./...

test:
	go test -v ./...

race:
	go test -race -v ./...

build-download_channels:
	@echo "download_channels..."
	go build -o build/download_channels cmd/download_channels/download_channels.go

build-download:
	@echo "Compiling download..."
	go build -o build/download cmd/download/download.go

build-search:
	@echo "Compiling search..."
	go build -o build/search cmd/search/search.go

build: build-download_channels build-download build-search

run-download_channels:
	go run cmd/download_channels/download_channels.go

run-download:
	go run cmd/download/download.go

run-search:
	go run cmd/search/search.go
