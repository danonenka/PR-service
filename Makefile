.PHONY: app-build app-run app-test app-clean app-docker-build app-docker-up app-docker-down app-migrate-up app-migrate-down app-logs app-status app-up app-down app-check-env

BINARY_NAME=pr-service
MAIN_PATH=./cmd/server
DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:5432/$(POSTGRES_DB)?sslmode=disable
DOCKER_DB_URL=postgres://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@postgres:5432/$(POSTGRES_DB)?sslmode=disable

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

app-build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

app-run:
	go run $(MAIN_PATH)

app-clean:
	go clean
	rm -f $(BINARY_NAME)

app-docker-build:
	docker-compose build

app-docker-up:
	docker-compose up -d

app-docker-down:
	docker-compose down

app-docker-logs:
	docker-compose logs -f

app-docker-status:
	docker-compose ps

app-migrate-up-local:
	migrate -path ./migrations -database "$(DB_URL)" up

app-migrate-down-local:
	migrate -path ./migrations -database "$(DB_URL)" down

app-migrate-up:
	docker-compose run --rm migrate -path /migrations -database "$(DOCKER_DB_URL)" up

app-migrate-down:
	docker-compose run --rm migrate -path /migrations -database "$(DOCKER_DB_URL)" down

app-up: app-docker-up app-migrate-up

app-down: app-docker-down

app-check-env:
	@which docker-compose > /dev/null || (echo "docker-compose not installed" && exit 1)
	@which migrate > /dev/null || (echo "migrate not installed" && exit 1)