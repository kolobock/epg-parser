.PHONY: build-darwin build-linux build-windows build lint test

build-darwin-amd:
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin-amd64/ ./...
build-darwin-arm:
	GOOS=darwin GOARCH=arm64 go build -o bin/darwin-arm64/ ./...
build-linux-amd:
	GOOS=linux GOARCH=amd64 go build -o bin/linux-amd64/ ./...
build-linux-arm:
	GOOS=linux GOARCH=arm64 go build -o bin/linux-arm64/ ./...
build-windows-amd:
	GOOS=windows GOARCH=amd64 go build -o bin/windows-amd64/ ./...
build-windows-arm:
	GOOS=windows GOARCH=arm64 go build -o bin/windows-arm64/ ./...
build: build-darwin-amd build-linux-arm
build-all: build-darwin-amd build-darwin-arm build-linux-amd build-linux-arm build-windows-amd build-windows-arm

test:
	go test -v -race ./...

default: build-darwin-amd

.ONESHELL:
GO_VER := $(shell go env GOVERSION | sed 's/go//' | cut -d. -f1,2)

LINT_VERSION_1.20 = v1.55
LINT_VERSION_1.21 = v1.56
LINT_VERSION_1.22 = v1.63.4
LINT_VERSION_1.23 = v1.64.8
LINT_VERSION_1.24 = v2.8.0

INSTALL_PATH_1.20 = github.com/golangci/golangci-lint/cmd/golangci-lint
INSTALL_PATH_1.21 = github.com/golangci/golangci-lint/cmd/golangci-lint
INSTALL_PATH_1.22 = github.com/golangci/golangci-lint/cmd/golangci-lint
INSTALL_PATH_1.23 = github.com/golangci/golangci-lint/cmd/golangci-lint

GOLINT_VERSION := $(or $(LINT_VERSION_$(GO_VER)),latest)
INSTALL_PATH := $(or $(INSTALL_PATH_$(GO_VER)),github.com/golangci/golangci-lint/v2/cmd/golangci-lint)

lint:
	@go install ${INSTALL_PATH}@${GOLINT_VERSION}
	golangci-lint run ./...

