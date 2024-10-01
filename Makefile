include .env
# Build the application
all: build

DB_URL=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432/${POSTGRES_DB}?sslmode=disable
DEV_URL=docker://postgres/15/dev
SCHEMA_FILE=file://database/schema.sql
MIGRATIONS_DIR=file://database/migrations

build:
	@echo "Building..."
	
	
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go



# Test the application
test:
	@echo "Testing..."
	@go test ./... -v



# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
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


generate:
	@sqlc generate

apply-schema:
	@echo "Applying schema to database..."
	atlas schema apply \
		--url "$(DB_URL)" \
		--to "$(SCHEMA_FILE)" \
		--dev-url "$(DEV_URL)"

migrate:
	@echo "Generating migration diff..."
	atlas migrate diff $(name) \
		--dir "$(MIGRATIONS_DIR)" \
		--to "$(SCHEMA_FILE)" \
		--dev-url "$(DEV_URL)"


.PHONY: all build run test clean watch generate
