# List of command names (directories)
COMMANDS := web cli

# Build all commands
build: $(COMMANDS)

# Build a specific command
$(COMMANDS):
	@echo "Building $@..."
	@go build -o bin/pm-$@ ./cmd/$@/main.go

# Clean all built commands
clean:
	@echo "Cleaning..."
	@rm -rf bin/*

# Run a specific command
run-%: build
	@echo "Running $*..."
	@./bin/$*

# Build the Docker container
build-docker:
	@echo "Building container..."
	$(eval GIT_TAG := $(shell git rev-parse --short HEAD))
	@echo  ----------------------------------------
	@echo   Git tag is: $(GIT_TAG), adding this to hex file now:
	@docker build -t plex-monitor:latest .
	@docker build -t plex-monitor:$(GIT_TAG) .

# Default target
.DEFAULT_GOAL := build
