APP_NAME=todo
BIN_DIR=bin

.PHONY: all build run clean fmt lint

# Default target
all: build

# Build binary into ./bin
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN_DIR)/$(APP_NAME) main.go

# Run with args (example: make run ARGS="create --name=test --tags=tag1,tag2")
run: build
	@$(BIN_DIR)/$(APP_NAME) $(ARGS)

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf $(BIN_DIR)

# Format code
fmt:
	go fmt ./...

# Lint (if you have golangci-lint installed)
lint:
	golangci-lint run
