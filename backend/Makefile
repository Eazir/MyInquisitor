.PHONY: all build run test clean sqlc migrate-up migrate-down

APP_NAME=myinquisitor
CMD_DIR=./cmd/api

all: build

build:
	go build -o bin/$(APP_NAME) $(CMD_DIR)

run:
	go run $(CMD_DIR)

test:
	go test ./... -v -cover

clean:
	rm -rf bin/

sqlc:
	sqlc generate

migrate-up:
	@echo "Run migrations manually or use a migration tool"

migrate-down:
	@echo "Rollback migrations manually or use a migration tool"

lint:
	golangci-lint run ./...

vet:
	go vet ./...
