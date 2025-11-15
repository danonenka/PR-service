.PHONY: build run test clean docker-build docker-up docker-down migrate-up migrate-down

BINARY_NAME=pr-service
MAIN_PATH=./cmd/server

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run:
	go run $(MAIN_PATH)

test:
	go test -v ./...

clean:
	go clean
	rm -f $(BINARY_NAME)

docker-build:
	docker-compose build

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_service?sslmode=disable" up

migrate-down:
	migrate -path ./migrations -database "postgres://postgres:postgres@localhost:5432/pr_service?sslmode=disable" down

up: docker-up

down: docker-down

