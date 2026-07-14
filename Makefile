include .env
export

MIGRATIONS_DIR := migrations

.PHONY: db-up
db-up:
	docker compose up -d --wait

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

.PHONY: test-integration
test-integration:
	CORDELL_INTEGRATION_TESTS=1 go test ./internal/infra/postgres -count=1 -v

.PHONY: fmt
fmt:
	gofmt -w .

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: run
run:
	go run ./cmd/cordell

.PHONY: check
check:
	make fmt
	make css-build
	make test
	make test-integration

.PHONY: css-build
css-build:
	npm run css:build

.PHONY: css-watch
css-watch:
	npm run css:watch

.PHONY: admin-create-operator
admin-create-operator:
	go run ./cmd/cordell-admin create-operator -username $(USERNAME)