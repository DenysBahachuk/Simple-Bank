run_postgres:
	sudo docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=admin -e POSTGRES_PASSWORD=adminpassword -d postgres:16.

create_db:
	sudo docker exec -it postgres16 createdb --username=admin --owner=admin Simple_Bank

drop_db:
	sudo docker exec -it postgres16 dropdb --username=admin Simple_Bank

migrate_up:
	migrate -path db/migration -database "postgresql://admin:adminpassword@localhost:5432/Simple_Bank?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgresql://admin:adminpassword@localhost:5432/Simple_Bank?sslmode=disable" -verbose down

sqlc_generate:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: run_postgres create_db drop_db migrate_up migrate_down sqlc_generate test