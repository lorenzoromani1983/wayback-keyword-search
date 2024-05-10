.PHONY: all build race build-downloader build-download build-search

all: lint test race build build-downloader build-download build-search

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

build-downloader:
	@echo "downloader..."
	go build -o build/downloader cmd/downloader/main.go

build-download:
	@echo "Compiling download..."
	go build -o build/download cmd/download/main.go

build-search:
	@echo "Compiling search..."
	go build -o build/search cmd/search/main.go

builds: build-downloader build-download build-search

run-downloader:
	go run cmd/downloader/main.go

run-download:
	go run cmd/download/main.go

run-search:
	go run cmd/search/main.go
