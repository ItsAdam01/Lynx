# Makefile for Lynx FIM

# Project variables
BINARY_NAME=lynx
VERSION=0.1.0
BUILD_DIR=bin
GO_FILES=$(shell find . -name "*.go" -not -path "./vendor/*")

# Standard targets
.PHONY: all build clean test run help

all: test build

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  build       Build the binary for the current OS/Arch"
	@echo "  test        Run all unit tests"
	@echo "  clean       Remove build artifacts"
	@echo "  release     Build cross-platform binaries (Linux amd64/arm64)"
	@echo "  help        Display this help message"

build: test
	@echo "Building Lynx FIM (current architecture)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) main.go

test:
	@echo "Running all tests..."
	go test -v ./...

clean:
	@echo "Cleaning build artifacts..."
	rm -rf $(BUILD_DIR)

# Cross-compilation for common platforms
.PHONY: release linux-amd64 linux-arm64

release: clean linux-amd64 linux-arm64
	@echo "Release builds complete in $(BUILD_DIR)/"

linux-amd64:
	@echo "Building for Linux (amd64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 main.go

linux-arm64:
	@echo "Building for Linux (arm64)..."
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 main.go
