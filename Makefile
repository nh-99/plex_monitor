# List of command names (directories)
COMMANDS := http cli discord
registryurl := localhost\:32000

# Build all commands
build: $(COMMANDS)

# Test all directories
test:
	@echo "Testing..."
	@go test ./...

# Build a specific command
$(COMMANDS):
	@echo "Building $@..."
	@go build --ldflags="-X 'plex_monitor/internal/buildflags.Version=$(shell git rev-parse --short HEAD)'" -o bin/pm-$@ ./cmd/$@/main.go

# Clean all built commands
clean:
	@echo "Cleaning..."
	@rm -rf bin/*

# Run a specific command
run-%: build
	@echo "Running $*..."
	@./bin/pm-$*

# Build the Docker container
build-docker:
	@echo "Building container..."
	$(eval GIT_TAG := $(shell git rev-parse --short HEAD))
	@echo  ----------------------------------------
	@echo   Git tag is: $(GIT_TAG), tagging container version
	@echo  ----------------------------------------
	@docker build -t $(registryurl)/plex-monitor:latest .
	@docker build -t $(registryurl)/plex-monitor:$(GIT_TAG) .
        @docker push $(registryurl)/plex-monitor:latest
        @docker push $(registryurl)/plex-monitor:$(GIT_TAG)

# Default target
.DEFAULT_GOAL := build
