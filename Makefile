postgres:
	docker run --name postgres_bank --network bank-network -p 5433:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine
createdb:
	docker exec -it postgres_bank createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres_bank dropdb simple_bank
migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgres://root:secret@localhost:5433/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate	
test:
	go test ./db/sqlc -v -cover

mock:
	mockgen -package mockdb -destination db/mock/store.go lesson/simple-bank/db/sqlc Store

.PHONY: postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 sqlc mock
