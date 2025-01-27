run_postgres:
	sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=adminpassword -d postgres:16.

create_db:
	sudo docker exec -it postgres16 createdb --username=admin --owner=admin simple_bank

drop_db:
	sudo docker exec -it postgres16 dropdb --username=admin simple_bank

migrate_up:
	migrate -path db/migration -database "postgresql://admin:adminpassword@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://admin:adminpassword@localhost:5432/simple_bank?sslmode=disable" -verbose down

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

.PHONY: run_postgres create_db drop_db migrate_up migrate_down sqlc_generate test server_run gen-mocks gen-docs