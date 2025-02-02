postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -d postgres:12-alpine
createdb:
	docker exec -it postgres createdb --username=root --owner=root bank-system

dropdb:
	docker exec -it postgres dropdb bank-system

migrateup:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:root@localhost:5432/bank-system?sslmode=disable" -verbose down

sqlc:
	sqlc generate
test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -package mockdb -destination db/mock/store.go github.com/spaghetti-lover/bank-system/db/sqlc Store
	
.PHONY: postgres createdb dropdb sqlc test server mock