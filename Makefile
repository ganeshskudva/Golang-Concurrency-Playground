# Makefile for building and running Go targets

# Variables
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
BINDIR=bin

# Directories
PUBSUB_DIR=cmd/pubsub
DEADLOCK_DIR=cmd/deadlockprevention
RATELIMITER_DIR=cmd/ratelimiter

LIB_DIR=pubsub
RATE_LIMITER_LIB=$(LIB_DIR)/ratelimiter/rate_limiter.go

# Source Files
PUBSUB_SRC=$(PUBSUB_DIR)/main.go
DEADLOCK_SRC=$(DEADLOCK_DIR)/main.go
RATELIMITER_SRC=$(RATELIMITER_DIR)/main.go

# Targets
PUBSUB_TARGET=$(BINDIR)/pubsub
DEADLOCK_PREVENTION_TARGET=$(BINDIR)/deadlock_prevention
RATELIMITER_TARGET=$(BINDIR)/rate_limiter

.PHONY: all clean run-pubsub run-deadlock run-ratelimiter

# Default target: Build all binaries
all: $(PUBSUB_TARGET) $(DEADLOCK_PREVENTION_TARGET) $(RATELIMITER_TARGET)

# Build pubsub
$(PUBSUB_TARGET): $(PUBSUB_SRC)
	mkdir -p $(BINDIR)
	$(GOBUILD) -o $(PUBSUB_TARGET) $(PUBSUB_SRC)

# Build deadlock prevention
$(DEADLOCK_PREVENTION_TARGET): $(DEADLOCK_SRC)
	mkdir -p $(BINDIR)
	$(GOBUILD) -o $(DEADLOCK_PREVENTION_TARGET) $(DEADLOCK_SRC)

# Build rate limiter
$(RATELIMITER_TARGET): $(RATELIMITER_SRC) $(RATE_LIMITER_LIB)
	mkdir -p $(BINDIR)
	$(GOBUILD) -o $(RATELIMITER_TARGET) $(RATELIMITER_SRC)

# Run pubsub
run-pubsub: $(PUBSUB_TARGET)
	$(PUBSUB_TARGET)

# Run deadlock prevention
run-deadlock: $(DEADLOCK_PREVENTION_TARGET)
	$(DEADLOCK_PREVENTION_TARGET)

# Run rate limiter
run-ratelimiter: $(RATELIMITER_TARGET)
	$(RATELIMITER_TARGET)

# Clean up build artifacts
clean:
	rm -rf $(BINDIR)

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
