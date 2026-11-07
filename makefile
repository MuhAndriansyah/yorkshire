.PHONY: run sqlc-generate oapi-bundle oapi-generate oapi-docs goose-up migration


SQLC_VERSION ?= v1.31.1
OAPI_CODEGEN_VERSION ?= v2.8.0
BUNDLED_OPENAPI := tmp/openapi.bundled.yaml
REDOCLY_VERSION ?= 2.39.0
GOOSE_VERSION ?= v3.15.0
DATABASE_URL := postgresql://postgres:postgres@localhost:5432/yorkshires?sslmode=disable
MIGRATION_NAME ?= $(filter-out $@,$(MAKECMDGOALS))

run:
	go run ./cmd/server/main.go

sqlc-generate:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@$(SQLC_VERSION) generate

oapi-bundle:
	npx --yes @redocly/cli@$(REDOCLY_VERSION) bundle contract/openapi.yaml -o $(BUNDLED_OPENAPI)

oapi-generate: oapi-bundle
	go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@$(OAPI_CODEGEN_VERSION) -config oapi-codegen.yaml $(BUNDLED_OPENAPI)

oapi-docs:
	npx --yes @redocly/cli@$(REDOCLY_VERSION) build-docs tmp/openapi.bundled.yaml -o tmp/docs.html --title "API Docs"
	open tmp/docs.html

goose-up:
	go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir migrations postgres "$(DATABASE_URL)" up


migration:
	@go run github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION) -dir migrations create $(MIGRATION_NAME) sql

# Cegah make error mencari target untuk argumen tambahan
%:
	@:
