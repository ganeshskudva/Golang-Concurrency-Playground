# Variables
BINARY_NAME := pubsub
BUILD_DIR := ./bin
SRC_DIR := ./cmd/pubsub
GO_FILES := $(shell find . -name '*.go')

# Default target
.PHONY: all
all: build

# Build target
.PHONY: build
build:
	@echo "Building the binary..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(SRC_DIR)/main.go

# Run the application
.PHONY: run
run: build
	@echo "Running the application..."
	@./$(BUILD_DIR)/$(BINARY_NAME)

# Clean target
.PHONY: clean
clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILD_DIR)

# Test target
.PHONY: test
test:
	@echo "Running tests..."
	go test ./...

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Running linter..."
	golangci-lint run

# Help target
.PHONY: help
help:
	@echo "Makefile Usage:"
	@echo "  make           Build the binary (default)"
	@echo "  make build     Build the binary"
	@echo "  make run       Run the application"
	@echo "  make clean     Remove build artifacts"
	@echo "  make test      Run all tests"
	@echo "  make fmt       Format code"
	@echo "  make lint      Run linter"
