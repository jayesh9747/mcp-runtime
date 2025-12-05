.PHONY: build test clean install dev fmt lint coverage

-include Makefile.runtime

# Binary name
BINARY_NAME=mcp-runtime
BUILD_DIR=bin

# Build the application for current platform
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/mcp-runtime

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	go clean

# Install dependencies
install:
	go mod download
	go mod tidy

# Development mode
dev: build
	./$(BUILD_DIR)/$(BINARY_NAME)

# Format code
fmt:
	go fmt ./...

# Lint code (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "golangci-lint not found. Please install from https://github.com/golangci/golangci-lint#install"; \
		exit 1; \
	fi

# Generate code coverage report
coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Build for Unix platforms (macOS and Ubuntu)
build-unix:
	@echo "Building for Unix platforms..."
	@mkdir -p $(BUILD_DIR)
	# macOS ARM64 (M1/M4)
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/mcp-runtime
	# macOS AMD64 (Intel)
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/mcp-runtime
	# Linux ARM64
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/mcp-runtime
	# Linux AMD64
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/mcp-runtime

# Install binary to system PATH
install-bin: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	@sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/
