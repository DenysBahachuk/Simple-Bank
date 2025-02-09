DB_URL=postgresql://admin:adminpassword@localhost:5432/simple_bank?sslmode=disable

run_postgres:
	sudo docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=adminpassword -d postgres:16.3

create_db:
	sudo docker exec -it postgres16 createdb --username=admin --owner=admin simple_bank

drop_db:
	sudo docker exec -it postgres16 dropdb --username=admin simple_bank

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

migrate_up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrate_down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc_generate:
	sqlc generate

test:
	go test -v -cover ./...

server_run:
	go run main.go

gen-mocks:
	mockgen -package mockdb -destination db/mock/store.go github.com/DenysBahachuk/Simple_Bank/db/sqlc Store

gen-docs:
	swag init -g server.go -d api,db/sqlc && swag fmt

db-docs:
	dbdocs build ./doc/db.dbml

db-schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

evans:
	evans --host localhost --port 9090 -r repl

.PHONY: run_postgres create_db drop_db migrate_up migrate_down sqlc_generate test server_run gen-mocks gen-docs new_migration db-docs db-schema proto evans