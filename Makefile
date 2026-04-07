dev:
	@docker compose up -d
	@$(MAKE) run

build:
	@go build  -o bin/gofollow ./cmd/api/

run: build
	@./bin/gofollow

test:
	@go test ./...

migrate-up:
	@migrate -path ./migrations -database "postgresql://admin:adminpassword@localhost:5432/gofollow?sslmode=disable" up

migrate-down:
	@migrate -path ./migrations -database "postgresql://admin:adminpassword@localhost:5432/gofollow?sslmode=disable" down

migrate-create:
	@migrate create -ext sql -dir ./migrations -seq $(name)

