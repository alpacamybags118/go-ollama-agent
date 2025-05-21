# Variables
BINARY_NAME=agent
BUILD_DIR=build
CMD_DIR=cmd/server

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin
GOFILES=$(wildcard *.go)

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

.PHONY: all build clean run test

## Build: builds the application
build:
	@echo "Building..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

## Install: installs the application
install:
	@echo "Installing..."
	@go install ./$(CMD_DIR)
	@echo "Install complete"

## Clean: cleans up binary files
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)
	@go clean
	@echo "Clean complete"

## Run: runs the application
run:
	@go run ./$(CMD_DIR)

## Test: runs go test with default values
test:
	@go test -v ./...

## Ollama: runs the Ollama container
ollama-up:
	@echo "Starting Ollama container..."
	@docker-compose up -d ollama
	@docker exec ollama ollama pull gemma3:1b

## Help: displays help
help:
	@echo "Make commands for go-ollama-agent:"
	@echo
	@grep -E '^##.*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = "## "}; {printf "\033[36m%-20s\033[0m %s\n", $$2, $$1}'
	@echo

# Default target
all: clean build