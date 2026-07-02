include .env
export

MIGRATIONS_DIR := migrations

.PHONY: db-up
db-up:
	docker compose up -d

.PHONY: db-down
db-down:
	docker compose down

.PHONY: db-logs
db-logs:
	docker compose logs -f postgres

.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATIONS_DIR) postgres "$(CORDELL_DATABASE_URL)" up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATIONS_DIR) postgres "$(CORDELL_DATABASE_URL)" down

.PHONY: migrate-status
migrate-status:
	goose -dir $(MIGRATIONS_DIR) postgres "$(CORDELL_DATABASE_URL)" status

.PHONY: test
test:
	go test ./...

.PHONY: fmt
fmt:
	gofmt -w .