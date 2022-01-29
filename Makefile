postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlcinit:
	sqlc init

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

server:
	go run main.go

sqlsession:
	docker exec -it postgres12 psql -U root -d simple_bank

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/tpmdigital/simplebank/db/sqlc Store
	
.PHONY: createdb dropdb migrateup migratedown migrateup1 migratedown1 sqlcinit sqlc test server sqlsession mock
