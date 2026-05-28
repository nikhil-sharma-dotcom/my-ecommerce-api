.PHONY: build run test clean migrate docker-build docker-up docker-down

APP_NAME=my-ecommerce-api
BUILD_DIR=./bin

build:
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)

migrate:
	psql $(DB_URL) -f migrations/001_create_users_table.sql
	psql $(DB_URL) -f migrations/002_create_products_table.sql
	psql $(DB_URL) -f migrations/003_create_orders_table.sql
	psql $(DB_URL) -f migrations/004_create_order_items_table.sql

docker-build:
	docker build -t $(APP_NAME):latest .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down -v

dev:
	air -c .air.toml
