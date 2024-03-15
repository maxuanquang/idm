database:
	docker run --name mysql -d -p 3306:3306 -e MYSQL_ROOT_PASSWORD=secret -e MYSQL_DATABASE=idm mysql:8.3.0 --character-set-server=utf8mb4 --collation-server=utf8mb4_unicode_ci

new_migration:
	migrate create -ext sql -dir ./internal/dataaccess/database/migrations/mysql -seq $(NAME)

up_migration:
	migrate -path ./internal/dataaccess/database/migrations/mysql -database "mysql://root:secret@tcp(0.0.0.0:3306)/idm?charset=utf8mb4&parseTime=True&loc=Local" -verbose up $(STEP)

down_migration:
	migrate -path ./internal/dataaccess/database/migrations/mysql -database "mysql://root:secret@tcp(0.0.0.0:3306)/idm?charset=utf8mb4&parseTime=True&loc=Local" -verbose down $(STEP)

proto:
	protoc \
	-I api \
	--go_out=./internal/generated \
	--go-grpc_out=./internal/generated \
	--validate_out="lang=go:./internal/generated" \
	--openapiv2_out=./api \
	--grpc-gateway_out ./internal/generated --grpc-gateway_opt generate_unbound_methods=true \
	api/idm.proto

wire:
	wire ./internal/wiring/wire.go

generate:
	make proto
	make wire

tidy:
	go mod tidy

.PHONY: proto, new_migration, up_migration, down_migration, tidy, wire, generate