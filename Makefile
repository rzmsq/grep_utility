# Project configuration
BINARY_NAME=grep_utility
MAIN_PATH=./cmd
PKG_PATH=./pkg/...

# Go configuration
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOFMT=gofmt
GOLINT=golangci-lint

# Build flags
BUILD_FLAGS=-ldflags="-s -w"
TEST_FLAGS=-v -race -coverprofile=coverage.out

# Default target
.PHONY: all
all: clean fmt lint test build

# Build the binary
.PHONY: build
build:
	$(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Build for different platforms
.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-linux-amd64 $(MAIN_PATH)

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)

.PHONY: build-all
build-all: build-linux build-windows build-mac

# Run tests
.PHONY: test
test:
	$(GOTEST) $(TEST_FLAGS) $(PKG_PATH)

.PHONY: test-short
test-short:
	$(GOTEST) -short $(PKG_PATH)

.PHONY: test-coverage
test-coverage: test
	$(GOCMD) tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Benchmarks
.PHONY: bench
bench:
	$(GOTEST) -bench=. -benchmem $(PKG_PATH)

# Format code
.PHONY: fmt
fmt:
	$(GOFMT) -s -w .

# Check formatting
.PHONY: fmt-check
fmt-check:
	@test -z "$$($(GOFMT) -l .)" || (echo "Code is not formatted. Run 'make fmt'"; exit 1)

# Lint code
.PHONY: lint
lint:
	$(GOLINT) run ./...

# Install linter if not present
.PHONY: install-linter
install-linter:
	@which $(GOLINT) > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin

# Run the application
.PHONY: run
run: build
	./$(BINARY_NAME)

# Run with example
.PHONY: run-example
run-example: build
	echo "hello world\ntest line\nhello again" | ./$(BINARY_NAME) hello

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*
	rm -f coverage.out coverage.html

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Verify dependencies
.PHONY: verify
verify:
	$(GOMOD) verify

# Install the binary
.PHONY: install
install:
	$(GOCMD) install $(BUILD_FLAGS) $(MAIN_PATH)

# Development workflow
.PHONY: dev
dev: fmt lint test build

# CI/CD pipeline
.PHONY: ci
ci: deps verify fmt-check lint test-coverage build

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Run clean, fmt, lint, test, build"
	@echo "  build         - Build the binary"
	@echo "  build-all     - Build for all platforms"
	@echo "  test          - Run tests with coverage"
	@echo "  test-short    - Run short tests"
	@echo "  test-coverage - Generate coverage report"
	@echo "  bench         - Run benchmarks"
	@echo "  fmt           - Format code"
	@echo "  fmt-check     - Check code formatting"
	@echo "  lint          - Run linter"
	@echo "  run           - Build and run the application"
	@echo "  run-example   - Run with example input"
	@echo "  clean         - Clean build artifacts"
	@echo "  tidy          - Tidy dependencies"
	@echo "  deps          - Download dependencies"
	@echo "  install       - Install the binary"
	@echo "  dev           - Development workflow"
	@echo "  ci            - CI/CD pipeline"
	@echo "  help          - Show this help"
