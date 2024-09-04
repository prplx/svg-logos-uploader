exist := $(wildcard .envrc)
ifneq ($(strip $(exist)),)
  include .envrc
endif

# Build the application
all: build

build:
	@echo "Building..."
	@templ generate
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

# Test the application
test:
	@echo "Testing..."
	@go test ./tests -v

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload
air:
	@if command -v air > /dev/null; then \
	    air; \
	    echo "Watching...";\
	else \
	    read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
	    if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	        echo "Watching...";\
	    else \
	        echo "You chose not to install air. Exiting..."; \
	        exit 1; \
	    fi; \
	fi

tailwindcss:
	@npx tailwindcss -i ./cmd/web/assets/css/input.css -o ./cmd/web/assets/css/output.css --watch

templ:
	@templ generate --watch cmd/web/**/*

.PHONY: all build run test clean run/watch air tailwindcss templ
