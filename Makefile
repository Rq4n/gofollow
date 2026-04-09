build:
	@go build  -o bin/gofollow ./cmd/api/

run: build
	@./bin/gofollow

test:
	@go test ./internal/service/...

migrate-create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

include .env

migrate-up:
	@migrate -path ./migrations -database "postgresql://$(POSTGRES_USER):$(POSTGRES_PASSWORD)@localhost:$(POSTGRES_PORT)/$(POSTGRES_DB)?sslmode=disable" up

migrate-down:
	@migrate -path ./migrations -database "postgresql://admin:adminpassword@localhost:5432/gofollow?sslmode=disable" down


