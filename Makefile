build:
	go build -o bin/urlshort main.go

lint:
	golangci-lint run

lint-fix:
	golangci-lint run --fix
	
format:
	golangci-lint fmt
	
test:
	go test -race ./...

app_up:
	docker compose --profile full up -d

app_down:
	docker compose --profile full down

db_up:
	docker compose up db -d

db_down:
	docker compose down db

db_connect:
	docker compose exec db psql -U postgres -d appdb

db_migrate:
	go tool \
	goose -dir ./db/migrations postgres "${DATABASE_URL}" \
	up
	
sqlc_generate:
	go tool sqlc generate

.PHONY: build lint lint-fix format test \
	app_up app_down \
	db_up db_down db_connect db_migrate \
	sqlc_generate
