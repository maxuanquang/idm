VERSION := 1.0.0
COMMIT_HASH := $(shell git rev-parse HEAD)
PROJECT_NAME := idm

.PHONY: database
database:
	docker run --name mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=idm mysql:8.3.0 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

.PHONY: new_migration
new_migration:
	migrate create -ext sql -dir ./internal/dataaccess/database/migrations/mysql -seq $(NAME)

.PHONY: up_migration
up_migration:
	migrate -path ./internal/dataaccess/database/migrations/mysql -database "mysql://root:secret@tcp(0.0.0.0:3306)/idm?charset=utf8mb4&parseTime=True&loc=Local" -verbose up $(STEP)

.PHONY: down_migration
down_migration:
	migrate -path ./internal/dataaccess/database/migrations/mysql -database "mysql://root:secret@tcp(0.0.0.0:3306)/idm?charset=utf8mb4&parseTime=True&loc=Local" -verbose down $(STEP)

.PHONY: proto
proto:
	protoc \
	-I api \
	--go_out=./internal/generated \
	--go-grpc_out=./internal/generated \
	--validate_out="lang=go:./internal/generated" \
	--openapiv2_out=./api/v1 \
	--grpc-gateway_out ./internal/generated --grpc-gateway_opt generate_unbound_methods=true \
	--experimental_allow_proto3_optional \
	api/idm.proto

.PHONY: wire
wire:
	wire ./internal/wiring/wire.go

.PHONY: openapi
openapi:
	openapi-generator-cli generate \
		-i api/v1/idm.swagger.json \
		-g javascript \
		-o web/src/api

.PHONY: generate
generate:
	make proto
	make wire
	make openapi

.PHONY: tidy
tidy:
	go mod tidy



.PHONY: build-web
build-web:
	cd web/src/api && npm install && cd ../../../
	cd web && npm install && npm run build && cd ..

.PHONY: build-backend
build-backend:
	go build \
		-ldflags "-X main.version=$(VERSION) -X main.commitHash=$(COMMIT_HASH)" \
		-o build/$(PROJECT_NAME) \
		cmd/$(PROJECT_NAME)/*.go

.PHONY: build
build:
	make build-web
	make build-backend

.PHONY: standalone-server
standalone-server:
	go run ./cmd/idm/main.go standalone-server

.PHONY: docker-compose-dev-up
docker-compose-dev-up:
	docker compose -f ./internal/deployments/docker-compose.dev.yml up -d

.PHONY: docker-compose-dev-down
docker-compose-dev-down:
	docker compose -f ./internal/deployments/docker-compose.dev.yml down

.PHONY: docker-compose-prod-up
docker-compose-prod-up:
	docker compose -f ./internal/deployments/docker-compose.prod.yml up -d

.PHONY: docker-compose-prod-down
docker-compose-prod-down:
	docker compose -f ./internal/deployments/docker-compose.prod.yml down

.PHONY: run
run:
	./build/idm standalone-server