# Variables
APP_NAME := gop
# VERSION := $(shell git describe --tags --abbrev=0)
# LDFLAGS := -ldflags "-X main.version=$(VERSION)"

# Default target
all: build

# Build the Go application
build:
	go build -o $(APP_NAME)

# Clean up
clean:
	rm -f $(APP_NAME)

# # Display the version (optional target)
# version:
# 	@echo $(VERSION)

.PHONY: all build clean
